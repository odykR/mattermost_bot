CREATE SCHEMA trello;
CREATE TABLE trello.user
(
    id            bigint PRIMARY KEY,
    login         text,
    password      text,
    tg_name       text,
    tg_username   text,
    register_time timestamp
);

CREATE TABLE trello.task
(
    id          SERIAL,
    user_id     bigint references trello.user (id),
    complexity  int,
    deadline    timestamp,
    description text
);

CREATE TABLE trello.team
(
    id   SERIAL PRIMARY KEY,
    name text
);

CREATE TABLE trello.user_team
(
    team_id int references trello.team (id),
    user_id bigint references trello.user (id)
);
