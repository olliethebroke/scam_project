package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// SelectLeaders считывает данные о лидерах из бд.
//
// В качестве выходных параметров возвращает мапу,
// содержащую лиги в качестве ключей и слайс указателей
// на тип данных Leader в качестве списка лидеров,
// и ошибку, если она возникла, в противном случае
// вместо неё будет возвращён nil.
func (pg *postgres) SelectLeaders() (map[int16][]*model.Leader, error) {
	// создаём sql запрос
	query, args, err := sq.Select("id", "position", "firstname", "blocks", "league").
		From("leaderboard").
		PlaceholderFormat(sq.Dollar).
		Limit(700).
		ToSql()
	if err != nil {
		return nil, err
	}

	// выполняя запрос, получаем несколько массивов по 100 элементов
	rows, err := pg.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	// создаём мапу лидеров в лигах
	leaders := make(map[int16][]*model.Leader)

	// проходимся по каждой строке, чтобы считать данные
	for rows.Next() {
		// считываем данные о лидере
		leader := &model.Leader{}
		var league int16
		err = rows.Scan(&leader.Id, &leader.Position, &leader.Firstname, &leader.Blocks, &league)
		if err != nil {
			return nil, err
		}

		// добавляем лидера в список лидеров конкретной лиги
		leaders[league] = append(leaders[league], leader)
	}

	// в случае успеха возвращаем результат
	return leaders, nil
}
