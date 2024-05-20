create table if not exists games (
    id integer primary key,
    name text not null,
    number text not null,
    player_uid integer not null,
    api_key text not null,
    start_time integer not null,
    tick_rate integer not null,
    production_rate integer not null,
    started integer not null check(started in (0, 1)),
    paused integer not null check(paused in (0, 1)),
    game_over integer not null check(game_over in (0, 1)),
    next_snapshot_at integer
);
