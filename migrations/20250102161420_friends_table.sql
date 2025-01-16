-- +goose Up
create table if not exists friends (
    id serial primary key,
    user_id bigint,
    friend_id bigint
);

-- +goose Down
drop table friends;
