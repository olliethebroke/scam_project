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

// NewPostgres инициализирует структуру postgres и
// возвращает указатель на неё.
//
// Входными параметрами метода являются контекст и
// dsn базы данных.
//
// Выходными параметрами метода являются указатель на
// тип postgres и ошибка, если она возникла, в противном
// случае будет возвращён nil.
func NewPostgres(ctx context.Context, dsn string) (*postgres, error) {
	pg := &postgres{}
	if err := pg.connect(ctx, dsn); err != nil {
		return nil, err
	}
	return pg, nil
}

// connect открывает соединение с бд, принимая дсн базы данных.
// Метод инициализирует структуру взаимодействия с бд
// Входными параметрами метода являются контекст и
// dsn базы данных.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае будет возвращён nil.
func (pg *postgres) connect(ctx context.Context, dsn string) error {
	var err error
	pg.pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	err = pg.pool.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Close закрывает соединение с бд
func (pg *postgres) Close() {
	pg.pool.Close()
}
