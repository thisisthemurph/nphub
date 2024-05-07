create table if not exists snapshots (
    id integer primary key,
    game_id integer references games(id),
    path text not null
);