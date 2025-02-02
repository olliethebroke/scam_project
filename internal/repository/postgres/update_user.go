package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// UpdateUser обновляет данные о пользователе в бд.
//
// Входными параметрами метода являются id пользователя
// и указатель на тип Update.
//
// Выходниым параметром метода является ошибка, если
// она возникла, в противном случае вместо неё будет
// возвращён nil.
func (pg *postgres) UpdateUser(id int64, update *model.Update) error {
	// создаём sql запрос
	query, args, err := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("blocks", update.Blocks).
		Set("record", update.Record).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	// обновляем нужные поля
	_, err = pg.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// возвращаем nil, если нет ошибок
	return nil
}
