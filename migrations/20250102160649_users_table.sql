-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id bigint not null,
    firstname text not null default 'You',
    blocks bigint not null default 0,
    record integer not null default 0,
    last_checkin timestamp not null default now(),
    days_streak smallint not null default 1,
    invited_friends smallint not null default 0,
    is_premium bool not null default false,
    created_at timestamp not null default now(),
    league smallint not null default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
