create table if not exists snapshots (
    id integer primary key,
    game_id integer,
    path text not null,
    created_at integer not null,
    constraint fk_games
        foreign key (game_id)
        references games(id)
        on delete cascade
);