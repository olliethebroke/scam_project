package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// InsertLeader добавляет лидера в базу данных.
//
// Входными параметрами метода являются контекст,
// указатель на тип Leader и лига лидера.
//
// Выходным параметром является ошибка,
// если она возникла, в противном случае
// вместо неё будет возвращён nil.
func (pg *postgres) InsertLeader(ctx context.Context, leader *model.Leader, league int16) error {
	// создаём sql запрос
	query, args, err := sq.Insert("leaderboard").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "position", "firstname", "blocks", "league").
		Values(leader.Id, leader.Position, leader.Firstname, leader.Blocks, league).
		ToSql()
	if err != nil {
		return err
	}

	// добавляем лидера в бд
	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// если не возникло ошибок
	// возвращаем nil
	return nil
}
