CREATE database if not exists insiderDb;

use insiderDb;

create table if not exists `team`
(
    id varchar(40) not null,
    name varchar(50) not null,
    homePoints json not null,
    awayPoints json not null,
    primary key (`id`)
);

create table if not exists `score`
(
    matchId varchar(40) not null,
    homeScore int not null,
    awayScore int not null,
    primary key (`matchId`)
);