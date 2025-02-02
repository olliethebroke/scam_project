package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
)

// SelectUserRole считывает информацию о роли пользователя.
//
// Метод на вход принимает id пользователя, чью роль необходимо
// считать.
//
// Выходными параметрами являются указатель на тип UserRole
// и ошибка, если она возникла,
// в противном случае вместо неё будет возращён nil.
func (pg *postgres) SelectUserRole(id int64) (*model.UserRole, error) {
	// формируем sql запрос
	query, args, err := sq.Select("role").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var role int16
	// считываем роль пользователя из бд
	err = pg.pool.QueryRow(context.Background(), query, args...).Scan(&role)
	if err != nil {
		return nil, err
	}
	// возвращаем данные о пользователе
	return &model.UserRole{
		Id:   id,
		Role: role,
	}, nil
}
