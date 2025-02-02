package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	"crypto_scam/pkg/hooks/telegram"
	sq "github.com/Masterminds/squirrel"
)

// InsertUser добавляет нового пользователя в базу даннх,
// возвращяя данные по умолчанию.
//
// На вход метод принимает id и имя нового пользователя.
//
// Выходными параметрами является указатель на тип User
// и ошибка, если она возникла, в противном случае
// вместо неё будет возвращён nil.
func (pg *postgres) InsertUser(id int64, firstname string) (*model.User, error) {
	// проверяем, есть ли пользователь в приватке
	isPremium, err := telegram.IfChatMember(int(id))
	if err != nil {
		return nil, err
	}

	// создаём sql запрос
	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id, firstname, is_premium").
		Values(id, firstname, isPremium).
		Suffix("RETURNING firstname").
		ToSql()
	if err != nil {
		return nil, err
	}

	// вносим пользователя в бд
	err = pg.pool.QueryRow(context.Background(), query, args...).Scan(&firstname)
	if err != nil {
		return nil, err
	}

	// возвращаем данные нового пользователя
	return &model.User{
		Id:             id,
		Firstname:      firstname,
		Blocks:         0,
		Record:         0,
		InvitedFriends: 0,
		DaysStreak:     1,
		IsPremium:      isPremium,
		League:         0,
		Award:          true,
	}, nil
}
