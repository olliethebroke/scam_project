-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule('*/10 * * * *', $$
    TRUNCATE TABLE leaderboard; -- Очищаем таблицу перед вставкой
    INSERT INTO leaderboard (id, position, username, blocks, league)
    SELECT id, position, username, blocks, league
    FROM (
         SELECT id, username, blocks, league,
                ROW_NUMBER() OVER (PARTITION BY league ORDER BY blocks DESC) as position
         FROM users
     ) AS ranked
    WHERE position <= 100
$$);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
