package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgres - структура, имплементирующая работу с PostgreSQL.
//
// pool - пул соедений с бд
type postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres возвращает указатель на структуру postgres
func NewPostgres() *postgres {
	return &postgres{}
}

// Connect открывает соединение с бд, принимая дсн базы данных.
// Метод инициализирует структуру взаимодействия с бд
func (pg *postgres) Connect(dsn string) error {
	ctx := context.Background()
	var err error
	pg.pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	return nil
}

// Close закрывает соединение с бд
func (pg *postgres) Close() {
	pg.pool.Close()
}
