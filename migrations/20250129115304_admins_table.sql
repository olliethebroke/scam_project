-- +goose Up
-- +goose StatementBegin
create table if not exists admins (
    id bigint,
    role smallint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table admins;
-- +goose StatementEnd
