### SETUP:
# get Go from https://golang.org/dl/
wget https://golang.org/dl/go1.15.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz

# Add /usr/local/go/bin to the PATH environment variable. You can do this by adding this line to your /etc/profile (for a system-wide installation) or $HOME/.profile:
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/{username}/go

# clone repo
cd /home/{username}/go
git clone https://github.com/Purely-Imaginary/referee-go.git 

# install dependencies
cd referee-go
go get ./...



## mysql CREATE TABLE script:

```mysql

create schema referee collate utf8mb4_0900_ai_ci;

create table referee.downloaded_url
(
	url varchar(300) null,
	match_id int null
);

create table referee.match_calculated
(
	id int auto_increment,
	time timestamp null,
	red_score int null,
	blue_score int null,
	red_avg float null,
	blue_avg float null,
	rating_change float null,
	raw_positions varchar(1000) null,
	constraint match_calculated_id_uindex
		unique (id)
);

alter table referee.match_calculated
	add primary key (id);

create table referee.player
(
	id int auto_increment,
	name varchar(50) null,
	wins int null,
	losses int null,
	goals_scored int null,
	goals_lost int null,
	win_rate int null,
	rating float null,
	constraint player_id_uindex
		unique (id)
);

alter table referee.player
	add primary key (id);

create table referee.player_snapshot
(
	id int auto_increment,
	player_id int null,
	player_name varchar(300) null,
	match_id int null,
	rating float null,
	is_red tinyint(1) null,
	constraint player_history_id_uindex
		unique (id)
);

alter table referee.player_snapshot
	add primary key (id);

create table referee.raw_match
(
	id int auto_increment,
	time timestamp null,
	positions varchar(300) null,
	red_score int null,
	blue_score int null,
	red_players varchar(300) null,
	blue_players varchar(300) null,
	constraint raw_match_id_uindex
		unique (id)
);

alter table referee.raw_match
	add primary key (id);


```
