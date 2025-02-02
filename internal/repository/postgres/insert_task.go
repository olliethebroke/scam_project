package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// InsertTask добавляет новое сформированное
// администратором задание в бд.
//
// В качестве входного параметра метод принимает
// указатель на тип данных Task.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (pg *postgres) InsertTask(task *model.Task) error {
	// создаём sql запрос
	query, args, err := sq.Insert("tasks").
		PlaceholderFormat(sq.Dollar).
		Columns("description", "reward", "action_type", "action_data").
		Values(task.Description, task.Reward, task.ActionType, task.ActionData).
		ToSql()
	if err != nil {
		return err
	}

	// добавляем задание в бд
	_, err = pg.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// возвращаем nil, если нет ошибок
	return nil
}
