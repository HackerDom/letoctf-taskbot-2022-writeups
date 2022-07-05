CREATE TABLE IF NOT EXISTS users
(
    id       TEXT primary key,
    username TEXT unique,
    pass     bytea
)