package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// SelectTasks забирает из бд список заданий.
//
// Входным параметром метода является контекст.
//
// Метод возвращает слайс указателей на тип Task
// и ошибку, если она возникла,
// в противном случае вместо неё будет возвращён nil.
func (pg *postgres) SelectTasks(ctx context.Context) ([]*model.Task, error) {
	// создаём sql запрос на получение списка заданий
	query, args, err := sq.Select("id", "description", "reward", "action_type", "action_data").
		PlaceholderFormat(sq.Dollar).
		From("tasks").
		ToSql()
	if err != nil {
		return nil, err
	}

	// получаем список всех заданий
	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// создаём слайс заданий
	tasks := make([]*model.Task, 0, 6)

	// проходимся по полученному списку и добавляем каждое задание в слайс
	for rows.Next() {
		// считываем информацию о задании
		t := &model.Task{}
		err = rows.Scan(&t.Id, &t.Description, &t.Reward, &t.ActionType, &t.ActionData)
		if err != nil {
			return nil, err
		}

		// добавляем задание в слайс
		tasks = append(tasks, t)
	}

	// возвращаем результат
	return tasks, nil
}
