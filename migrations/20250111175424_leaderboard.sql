-- +goose Up
-- +goose StatementBegin
create table leaderboard (
    id bigint not null,
    position smallint not null default 100,
    firstname text not null default 'You',
    blocks bigint not null default 0,
    league smallint not null default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table leaderboard;
-- +goose StatementEnd
