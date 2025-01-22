-- +goose Up
create table if not exists friends (
    id serial primary key,
    invited_user_id bigint,
    inviting_user_id bigint
);

-- +goose Down
drop table friends;
