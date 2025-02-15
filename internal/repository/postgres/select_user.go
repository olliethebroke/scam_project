package postgres

import (
	"context"
	"crypto_scam/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
	"time"
)

// SelectUser считывает данные о пользователе из бд.
//
// В качестве параметра принимает контекст и id пользователя,
// данные о котором нужно считать.
//
// Выходными параметрами метода являются указатель
// на тип User и ошибка, если она возникла,
// в противном случае вместо неё будет возращён nil.
func (pg *postgres) SelectUser(ctx context.Context, id int64) (*model.User, error) {
	// создаём sql запрос
	query, args, err := sq.Select("id", "firstname", "blocks", "record", "last_checkin", "days_streak", "invited_friends", "is_premium", "league").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	// получаем данные о пользователе
	var info = &model.User{}
	var lastCheckin = &time.Time{}
	err = pg.pool.QueryRow(ctx, query, args...).Scan(
		&info.Id, &info.Firstname,
		&info.Blocks, &info.Record, lastCheckin,
		&info.DaysStreak, &info.InvitedFriends,
		&info.IsPremium, &info.League)
	if err != nil {
		return nil, err
	}

	// обновляем серию ежедневных заходов
	award, err := pg.updateStreak(ctx, id, lastCheckin, &info.DaysStreak)
	if err != nil {
		return nil, err
	}

	// присваиваем переменной award в респонсе, нужно ли выдавать игроку награду за вход
	info.Award = award

	// возвращаем результат
	return info, nil
}

// updateStreak обновляет стрик ежедневных заходов в игру.
//
// Входными параметрами метода являются контекст,
// идентификатор пользователя, которому нужно обновить стрик,
// послденее время сбора награды, в виде указателя на Time
// и указатель на int16, отвечающий за обновление стрика.
//
// Выходными параметрами метода являются булевый тип данных
// и ошибка, если она возникла, в противом случае будет возвращён nil.
func (pg *postgres) updateStreak(ctx context.Context, id int64, lastCheckin *time.Time, streak *int16) (bool, error) {
	// время проверки
	now := time.Now()
	// разность текущего времени и времени получения последней daily награды
	diff := now.Sub(*lastCheckin)
	// если зашёл на следующий день
	if diff.Hours() >= 24 && diff.Hours() < 48 {
		// увеличиваем стрик на 1
		*streak += 1
		// создаём sql запрос
		query, args, err := sq.Update("users").
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{"id": id}).
			Set("days_streak", *streak).
			Set("last_checkin", now).
			ToSql()
		if err != nil {
			return false, err
		}
		// увеличиваем стрик и обновляем время захода
		_, err = pg.pool.Exec(ctx, query, args...)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// если игрок не заходил более двух суток
	if diff.Hours() >= 48 {
		// обнуляем стрик до 1
		*streak = 1
		// создаём sql запрос
		query, args, err := sq.Update("users").
			PlaceholderFormat(sq.Dollar).
			Set("days_streak", *streak).
			Set("last_checkin", now).
			Where(sq.Eq{"id": id}).
			ToSql()
		if err != nil {
			return false, err
		}
		// обнуляем стрик
		_, err = pg.pool.Exec(ctx, query, args...)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// если ни одно из условий не выполняется - игрок зашёл в течение суток; награды нет
	return false, nil
}
