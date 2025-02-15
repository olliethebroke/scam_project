package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

// InsertFriendship создаёт новую дружбу в бд
// и увеличивает счётчик друзей у приглашающего пользователя.
//
// Входными параметрами метода являются контекст,
// идентификатор приглашённого пользователя и
// идентификатор приглашающего пользователя.
//
// Метод возвращает ошибку, если она возникла,
// в противном случае вместо неё будет возвращён nil.
func (pg *postgres) InsertFriendship(ctx context.Context, invitedUserId int64, invitingUserId int64) error {
	// создаём sql запрос
	query, args, err := sq.Insert("friends").
		PlaceholderFormat(sq.Dollar).
		Columns("invited_user_id", "inviting_user_id").
		Values(invitedUserId, invitingUserId).
		ToSql()
	if err != nil {
		return err
	}

	// создаём запись дружбы в таблице friends
	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// создаём sql запрос
	query, args, err = sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("invited_friends", sq.Expr("invited_friends + 1")).
		Where(sq.Eq{"id": invitingUserId}).
		ToSql()
	if err != nil {
		return err
	}

	// увеличиваем счётчик друзей у приглашающего на единицу
	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// возвращем nil, если нет ошибок
	return nil
}
