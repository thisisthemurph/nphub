create table if not exists games (
    id integer primary key,
    number text not null,
    api_key text not null
);
