-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    ID serial primary key
);

create table if not exists events
    (
        ID serial primary key,
        UserID int not null,
        Title varchar(20) not null,
        Description text not null,
        EndEvent timestamp not null,
        StartEvent timestamp not null
);

-- +goose Down
-- +goose StatementBegin
drop table events;
drop table users;
-- +goose StatementEnd