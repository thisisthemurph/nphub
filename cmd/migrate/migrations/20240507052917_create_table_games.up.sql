create table if not exists games (
    id integer primary key,
    number text not null,
    player_uid integer not null,
    api_key text not null,
    start_time integer not null,
    tick_rate integer not null
);
