import replay, os, json
import threading

matches = {}
threads = []

def threadedAnalysis(path):
    with open(path, "rb") as game:
            print(path + ' starting...', end="")
            bin = replay.Replay(game.read()) 
            time = path[-26:-10].replace('h',':')
            time = time[:10] + ' ' + time[11:]
            meaningfulStates = []
            actualPlayers = {}
            matchLengthInTicks = 0
            for tick in bin[1]:
                if tick.state.value > 2:
                    for player in tick.players:
                        if player.id not in actualPlayers:
                            actualPlayers[player.id] = {'team': player.team.name, 'ticks': 0}
                        actualPlayers[player.id]['team'] = player.team.name
                        actualPlayers[player.id]['ticks'] = actualPlayers[player.id]['ticks'] + 1
                    matchLengthInTicks += 1
                    meaningfulStates.append(tick)

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
                'time': time,
                'teams': teams,
                'score': {
                    'Red': lastState[0].score[0],
                    'Blue': lastState[0].score[1]
                },
                'rawPositionsAtEnd': rawPositions
            }

            print(' processed')
            s = json.dumps(saveData, default=lambda x: x.__dict__)
            with open('../files/replayData/' + path[13:] + '.json', 'w+') as f:
                f.write(s)

for subdir, dirs, files in os.walk('preprocessed/'):
        for file in files:
            if file.split('.')[-1] != "bin":
                continue
            path = os.path.join(subdir, file)
            threadedAnalysis(path)
            os.remove(path)