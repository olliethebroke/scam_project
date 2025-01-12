package postgres

import (
	"context"
	"crypto_scam/pkg/hooks/telegram"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	Info      UserInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
}
type UserInfo struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	Blocks         int64  `json:"blocks"`
	Record         int    `json:"record"`
	DaysStreak     int    `json:"days_streak"`
	InvitedFriends int    `json:"invited_friends"`
	IsPremium      bool   `json:"is_premium"`
	League         int16  `json:"league"`
}
type Leader struct {
	Id       int64  `json:"id"`
	Position int16  `json:"position"`
	Username string `json:"username"`
	Blocks   int64  `json:"blocks"`
}
type UpdateRequest struct {
	Blocks         int64 `json:"blocks"`
	Record         int   `json:"record"`
	InvitedFriends int   `json:"invited_friends"`
	IsPremium      bool  `json:"is_premium"`
}
type CreateRequest struct {
	Username string `json:"username"`
}

var pool *pgxpool.Pool

// Connect открывает соединение с бд
func Connect(dsn string) error {
	ctx := context.Background()
	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	return nil
}

// InsertUser записывает пользователя в бд
func InsertUser(id int64, username string) (*UserInfo, error) {
	isPremium := true
	if err := telegram.IfChatMember(int(id)); err != nil {
		isPremium = false
	}
	// создаём sql запрос
	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id, username, is_premium").
		Values(id, username, isPremium).
		Suffix("RETURNING id, username").
		ToSql()

	if err != nil {
		return nil, err
	}
	info := UserInfo{}
	err = pool.QueryRow(context.Background(), query, args...).Scan(&info.Id, &info.Username)
	if err != nil {
		return nil, err
	}
	return &UserInfo{
		Id:             info.Id,
		Username:       info.Username,
		Blocks:         0,
		Record:         0,
		InvitedFriends: 0,
		DaysStreak:     1,
		IsPremium:      isPremium,
		League:         0,
	}, nil
}

// SelectUser считывает данные о пользователе из бд
func SelectUser(id int64) (*UserInfo, error) {
	// создаём sql запрос
	query, args, err := sq.Select("id", "username", "blocks", "record", "last_checkin", "days_streak", "invited_friends", "is_premium", "league").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var info = &UserInfo{}
	err = pool.QueryRow(context.Background(), query, args...).Scan(
		&info.Id, &info.Username,
		&info.Blocks, &info.Record,
		&info.DaysStreak, &info.InvitedFriends,
		&info.IsPremium, &info.League)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// SelectLeaders считывает данные о лидерах из бд
func SelectLeaders() (map[string][]*Leader, error) {
	// создаём sql запрос
	query, args, err := sq.Select("id", "position", "username", "blocks", "league").
		From("leaderboard").
		PlaceholderFormat(sq.Dollar).
		Limit(700).
		ToSql()
	if err != nil {
		return nil, err
	}
	// выполняя запрос, получаем несколько массивов по 100 элементов
	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	leaders := make(map[string][]*Leader)
	// проходимся по каждой строке, чтобы считать данные
	for rows.Next() {
		leader := &Leader{}
		var league string
		err = rows.Scan(&leader.Id, &leader.Position, &leader.Username, &leader.Blocks, &league)
		if err != nil {
			return nil, err
		}
		leaders[league] = append(leaders[league], leader)
	}
	return leaders, nil
}

// UpdateUser обновляет инфу о пользователе в бд
func UpdateUser(id int64, info *UpdateRequest) error {
	// создаём sql запрос
	query, args, err := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("blocks", info.Blocks).
		Set("record", info.Record).
		Set("invited_friends", info.InvitedFriends).
		Set("is_premium", info.IsPremium).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Close закрывает соединение с бд
func Close() {
	pool.Close()
}
