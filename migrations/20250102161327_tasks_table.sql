-- +goose Up
create table if not exists tasks (
    id serial primary key,
    description text,
    reward integer not null default 0,
    action_type text,
    action_data text
);

-- +goose Down
drop table tasks;
