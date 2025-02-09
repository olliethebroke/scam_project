package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

// InsertCompletedTask добавляет выполненное пользователем задание
// в базу данных.
//
// Входными параметрами метода являются идентификатор пользователя,
// выполнившего задание, и идентификатор самого задания.
//
// Выходным параметром метода является ошибка, если она не возникла,
// то вместо неё будет возвращён nil.
func (pg *postgres) InsertCompletedTask(userId int64, taskId int16) error {
	// создаём sql запрос
	query, args, err := sq.Insert("completed_tasks").
		PlaceholderFormat(sq.Dollar).
		Columns("task_id", "user_id").
		Values(taskId, userId).
		ToSql()
	if err != nil {
		return err
	}

	// добавляем выполненное задание в бд
	_, err = pg.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// если ошибок не возникло
	// возвращаем nil
	return nil
}
