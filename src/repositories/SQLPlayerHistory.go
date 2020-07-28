package repositories

type SQLPlayerHistory struct {
	ID       int64   `db:"uid,pk"`
	MatchID  int64   `db:"match_id"`
	PlayerID int64   `db:"player_id"`
	Time     int64   `db:"timestamp"`
	Rating   float32 `db:"rating"`
	IsRed    bool    `db:"is_red"`
}

func (u *SQLPlayerHistory) TableName() string {
	return "player_history"
}

/*
create table player_history
(
	id int not null,
	player_id int null,
	time int null,
	rating float null,
	is_red bool null
);

create unique index player_history_id_uindex
	on player_history (id);

alter table player_history
	add constraint player_history_pk
		primary key (id);

*/
