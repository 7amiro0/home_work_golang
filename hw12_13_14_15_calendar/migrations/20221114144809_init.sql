-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table if not exists users
(
    id   serial primary key,
    name varchar(20) unique not null
);

create table if not exists events
(
    id serial primary key,
    userID int not null,
    userName varchar(20),
    title varchar(20) not null,
    description text not null,
    notify timestamp not null,
    startEvent timestamp not null,
    endEvent timestamp not null,
    constraint userID foreign key(userID) references users(id),
    constraint userName foreign key(userName) references users(name)
);

create or replace function new_event(
    userName varchar(20),
    title varchar(20),
    description text,
    notify timestamp,
    startEvent timestamp,
    endEvent timestamp
) returns integer as $$
declare
    identifier integer := 0;
    eventID integer := 0;
begin

    select id from users where name = $1 into identifier;

    if identifier isnull or identifier = 0 then
        insert into users (name) values ($1) returning users.id into identifier;
    end if;

    insert into
        events (userID, userName, title, description, notify, startEvent, endEvent)
    values
        (identifier, $1, $2, $3, $4, $5, $6)
    returning
        events.id into eventID;

    return eventID;
end; $$
    language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table events;
drop table users;
-- +goose StatementEnd
