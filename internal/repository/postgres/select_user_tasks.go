package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
	"slices"
)

// SelectUserTasks считывает из бд массив заданий пользователя по его id.
//
// Входными параметрами метода является id пользовталя,
// задания которого необходимо считать.
//
// Выходными параметрами метода являются слайс указателей
// на тип User и ошибка, если она возникла,
// в противном случае вместо неё будет возращён nil.
func (pg *postgres) SelectUserTasks(id int64) ([]*model.Task, error) {
	// создаём sql запрос на получение выполненных пользователем заданий
	query, args, err := sq.Select("task_id").
		PlaceholderFormat(sq.Dollar).
		From("completed_tasks").
		Where(sq.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	// получаем список выполненных заданий пользователя
	rows, err := pg.pool.Query(context.Background(), query, args...)

	// создаём слайс выполненных заданий
	completedTasks := make([]int16, 0, 6)
	if err != nil {
		return nil, err
	}

	// отмечаем выполненные пользователем задания в слайсе completedTasks
	for rows.Next() {
		// считываем id выполненного задания
		var completedTaskId int16
		err = rows.Scan(&completedTaskId)
		if err != nil {
			return nil, err
		}
		// добавляем id задания в слайс выполненных заданий
		completedTasks = append(completedTasks, completedTaskId)
	}

	// создаём sql запрос на получение списка заданий
	query, args, err = sq.Select("id", "description", "reward", "action_type", "action_data").
		PlaceholderFormat(sq.Dollar).
		From("tasks").
		ToSql()
	if err != nil {
		return nil, err
	}

	// получаем список всех заданий
	rows, err = pg.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	// создаём слайс заданий
	tasks := make([]*model.Task, 0, 6)

	// проходимся по полученному списку и добавляем каждое задание в слайс
	for rows.Next() {
		t := &model.Task{}
		// считываем информацию о задании
		err = rows.Scan(&t.Id, &t.Description, &t.Reward, &t.ActionType, &t.ActionData)
		if err != nil {
			return nil, err
		}

		// если задание есть в списке выполненных, то ставим флаг
		t.IsCompleted = slices.Contains(completedTasks, t.Id)

		// добавляем задание в слайс
		tasks = append(tasks, t)
	}
	// возвращаем результат
	return tasks, nil
}
