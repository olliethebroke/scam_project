package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// SelectTask считывает данные о задании из бд.
//
// В качестве параметра принимает контекст и id задания,
// данные о котором нужно считать.
//
// Выходными параметрами метода являются указатель
// на тип Task и ошибка, если она возникла,
// в противном случае вместо неё будет возращён nil.
func (pg *postgres) SelectTask(ctx context.Context, id int16) (*model.Task, error) {
	// создаём sql запрос
	query, args, err := sq.Select("description", "reward", "action_type", "action_data").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		From("tasks").
		ToSql()
	if err != nil {
		return nil, err
	}

	// считываем задание по id
	task := &model.Task{}
	err = pg.pool.QueryRow(ctx, query, args...).Scan(
		&task.Description, &task.Reward,
		&task.ActionType, &task.ActionData)
	if err != nil {
		return nil, err
	}

	// возвращаем результат
	return task, nil
}
