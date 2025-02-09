package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

// InsertAdmin вносит нового админа в бд.
//
// Входным параметром является id нового администратора.
//
// Метод возвращает ошибку, если она возникла,
// в противном случае вместо неё будет возвращён nil.
func (pg *postgres) InsertAdmin(id int64) error {
	// создаём sql запрос
	query, args, err := sq.Insert("admins").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "role").
		Values(id, 1).
		ToSql()
	if err != nil {
		return err
	}

	// вносим админа в бд
	_, err = pg.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// если ошибок не возникло
	// возвращаем nil
	return nil
}
