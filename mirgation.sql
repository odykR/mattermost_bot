CREATE SCHEMA bot;
CREATE TABLE bot.user
(
    id            bigint PRIMARY KEY,
    login         text,
    password      text,
    tg_name       text,
    tg_username   text,
    register_time timestamp
);

CREATE TABLE bot.task
(
    id          SERIAL,
    user_id     bigint references bot.user (id),
    complexity  int,
    deadline    timestamp,
    description text
);

CREATE TABLE bot.team
(
    id   SERIAL PRIMARY KEY,
    name text
);

CREATE TABLE bot.user_team
(
    team_id int references bot.team (id),
    user_id bigint references bot.user (id)
);
