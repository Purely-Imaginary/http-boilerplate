import replay, os, json, math
import threading

matches = {}
threads = []

def calculateDistance(x1,y1,x2,y2):
    dist = math.sqrt((x2 - x1)**2 + (y2 - y1)**2)
    return dist

def hasVectorChanged(ball1, ball2):
    vector1 = math.atan2(ball1.vx, ball1.vy) * (180/math.pi)
    vector2 = math.atan2(ball2.vx, ball2.vy) * (180/math.pi)
    change = abs(vector1 - vector2)
    return change > 0.5

def threadedAnalysis(path):
    with open(path, "rb") as game:
            print(path + ' starting...', end="")
            bin = replay.Replay(game.read()) 
            time = path[-26:-10].replace('h',':')
            time = time[:10] + ' ' + time[11:]
            meaningfulStates = []
            actualPlayers = {}
            matchLengthInTicks = 0
            tickCounter = 0
            goals = []
            for tick in bin[1]:
                if tick.state.value == 3:
                    for player in tick.players:
                        if player.id not in actualPlayers:
                            actualPlayers[player.id] = {'team': player.team.name, 'ticks': 0}
                        actualPlayers[player.id]['team'] = player.team.name
                        actualPlayers[player.id]['ticks'] = actualPlayers[player.id]['ticks'] + 1
                    matchLengthInTicks += 1
                    meaningfulStates.append(tick)
                
                if (tick.state.value == 4 and bin[1][tickCounter-1].state.value == 3 and tickCounter != 0):
                    goalScorerId = -1
                    goalShotTime = 0
                    # TODO: goalShotSpeed
                    goalSide = "Red"

                    if len(goals) == 8:
                        a=0

                    if (bin[1][tickCounter-1].score[0] == tick.score[0]):
                        goalSide = "Blue"
                    for i in reversed(range(tickCounter-1)):
                        if goalScorerId != -1:
                            break
                        targetTick = bin[1][i]
                        distances = []
                        for player in targetTick.players:
                            if player.team.name != goalSide:    
                                continue
                            distanceFromBall = calculateDistance(targetTick.ball.x, targetTick.ball.y, player.disc.x, player.disc.y)
                            distances.append({"id": player.id, "dist":distanceFromBall})

                            if distanceFromBall < 30 and hasVectorChanged(targetTick.ball, bin[1][i+1].ball):
                                goalScorerId = player.id
                                goalShotTime = round(targetTick.gameTime, 3)
                                break
                        
                        a=0

                    goals.append({
                        "goalTime": round(tick.gameTime, 3),
                        "goalScorerId": goalScorerId,
                        "goalShotTime": goalShotTime,
                        "goalSide": goalSide,
                        "goalTravelTime": round(tick.gameTime - goalShotTime, 3)
                    })

                tickCounter+=1
            if meaningfulStates.__len__() == 0:
                print(' empty')
                return

            lastState = meaningfulStates[-1:]
            rawPositions = ""
            teams = {'Red': [], 'Blue': []}
            for player in lastState[0].players:
                rawPositions += str(player.disc.x) + "-" + str(player.disc.y) + "|"

            for playerId in actualPlayers:
                if (actualPlayers[playerId]['ticks'] / matchLengthInTicks > 0.6):
                    teams[actualPlayers[playerId]['team']].append(bin[0][playerId])                
             
            saveData = {
                'gameTime': round(lastState[0].gameTime, 3),
                'time': time,
                'teams': teams,
                'goalsData': goals,
                'score': {
                    'Red': lastState[0].score[0],
                    'Blue': lastState[0].score[1]
                },
                'rawPositionsAtEnd': rawPositions
            }

            print(' processed')
            s = json.dumps(saveData, default=lambda x: x.__dict__, sort_keys=True, indent=4)
            with open('../files/replayData/' + path[13:] + '.json', 'w+') as f:
                f.write(s)

for subdir, dirs, files in os.walk('preprocessed/'):
        for file in files:
            if file.split('.')[-1] != "bin":
                continue
            path = os.path.join(subdir, file)
            threadedAnalysis(path)
            os.remove(path)