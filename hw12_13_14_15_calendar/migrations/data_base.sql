-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    ID serial primary key,
    Title text,
    IDUser int,
    Description text,
    StartEvent timestamp,
    EndEvent timestamp
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd