package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

// DeleteUser удаляет пользователя из базы данных.
//
// Метод принимает id пользователя,
// которого требуется удалить.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (pg *postgres) DeleteUser(id int64) error {
	// создаём sql запрос на удаление пользователя
	query, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	// удаляем пользователя из бд
	_, err = pg.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// если не возникло ошибок
	// возвращаем nil
	return nil
}
