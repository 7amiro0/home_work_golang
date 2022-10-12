-- +goose Up
-- +goose StatementBegin
create table events
(
    ID serial primary key,
    IDUser int,
    Title varchar(20),
    Description text,
    EndEvent timestamp,
    StartEvent timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events
-- +goose StatementEnd
