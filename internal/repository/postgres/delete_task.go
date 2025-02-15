package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

// DeleteTask удаляет задание из базы данных.
//
// Метод принимает контекст и
// id задания, которое нужно удалить.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (pg *postgres) DeleteTask(ctx context.Context, id int16) error {
	// создаём sql запрос на удаление нужного задания из бд
	query, args, err := sq.Delete("tasks").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	// удаляем задание
	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// если всё прошло без ошибок
	// возвращаем nil
	return nil
}
