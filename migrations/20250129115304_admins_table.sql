-- +goose Up
-- +goose StatementBegin
create table if not exists admins (
    admin_id bigint,
    role smallint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table admins;
-- +goose StatementEnd
