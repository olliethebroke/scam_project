package postgres

import (
	"context"
	"crypto_scam/pkg/hooks/telegram"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Leader содержит информацию для отображения лидера лиги
type Leader struct {
	Id       int64  `json:"id"`
	Position int16  `json:"position"`
	Username string `json:"username"`
	Blocks   int64  `json:"blocks"`
}

// CreateUserRequest содержит поле запроса от клиента на создание пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}

type SelectUserResponse struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	Blocks         int64  `json:"blocks"`
	Record         int    `json:"record"`
	DaysStreak     int    `json:"days_streak"`
	InvitedFriends int    `json:"invited_friends"`
	IsPremium      bool   `json:"is_premium"`
	League         int16  `json:"league"`
	Award          bool   `json:"award"`
}

// UpdateUserRequest содержит поля запроса от клиента на обновление данных о пользователе
type UpdateUserRequest struct {
	Blocks         int64 `json:"blocks"`
	Record         int   `json:"record"`
	InvitedFriends int   `json:"invited_friends"`
	IsPremium      bool  `json:"is_premium"`
}

// переменная - пул соедений с бд
var pool *pgxpool.Pool

// Connect открывает соединение с бд, принимая дсн базы данных
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
// возвращает ошибку и данные пользователя
func InsertUser(id int64, username string) (*SelectUserResponse, error) {
	// задаём дефолтный флаг премиум пользователя
	isPremium := true
	// проверяем, есть ли пользователь в приватке
	if err := telegram.IfChatMember(int(id)); err != nil {
		// если нет - убираем флаг
		isPremium = false
	}
	// создаём sql запрос
	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id, username, is_premium").
		Values(id, username, isPremium).
		Suffix("RETURNING username").
		ToSql()
	if err != nil {
		return nil, err
	}
	// вносим пользователя в бд
	err = pool.QueryRow(context.Background(), query, args...).Scan(&username)
	if err != nil {
		return nil, err
	}
	// возвращаем данные нового пользователя
	return &SelectUserResponse{
		Id:             id,
		Username:       username,
		Blocks:         0,
		Record:         0,
		InvitedFriends: 0,
		DaysStreak:     1,
		IsPremium:      isPremium,
		League:         0,
		Award:          true,
	}, nil
}

// SelectUser считывает данные о пользователе из бд
// возвращает ошибку и данные пользователя
func SelectUser(id int64) (*SelectUserResponse, error) {
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
	var info = &SelectUserResponse{}
	var lastCheckin = &time.Time{}
	// получаем данные о пользователе
	err = pool.QueryRow(context.Background(), query, args...).Scan(
		&info.Id, &info.Username,
		&info.Blocks, &info.Record, lastCheckin,
		&info.DaysStreak, &info.InvitedFriends,
		&info.IsPremium, &info.League)
	if err != nil {
		return nil, err
	}
	// обновляем серию ежедневных заходов
	award, err := updateStreak(id, lastCheckin, info.DaysStreak)
	if err != nil {
		return nil, err
	}
	// присваиваем переменной award в респонсе, нужно ли выдавать игроку награду за вход
	info.Award = award
	return info, nil
}

// SelectLeaders считывает данные о лидерах из бд;
// возвращает ошибку и мапу, содержащую для каждой лиги(int16) слайс из лидеров
func SelectLeaders() (map[int16][]*Leader, error) {
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
	// создаём мапу лидеров в лигах
	leaders := make(map[int16][]*Leader)
	// проходимся по каждой строке, чтобы считать данные
	for rows.Next() {
		leader := &Leader{}
		var league int16
		// считываем данные о лидере
		err = rows.Scan(&leader.Id, &leader.Position, &leader.Username, &leader.Blocks, &league)
		if err != nil {
			return nil, err
		}
		// добавляем лидера в список лидеров конкретной лиги
		leaders[league] = append(leaders[league], leader)
	}
	return leaders, nil
}

// UpdateUser обновляет инфу о пользователе в бд;
// аргумент UpdateUserRequest - полученная информация о пользователе
func UpdateUser(id int64, info *UpdateUserRequest) error {
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
	// обновляем нужные поля
	_, err = pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return nil
}

// updateStreak обновляет стрик ежедневных заходов в игру;
// аргумент lastCheckin - время последнего захода, за который давалась daily награда;
// аргумент streak - текущий стрик игрока с указанным id;
// возвращает ошибку и флаг о том, нужна ли награда
func updateStreak(id int64, lastCheckin *time.Time, streak int) (bool, error) {
	// время проверки
	now := time.Now()
	// разность текущего времени и времени получения последней daily награды
	diff := now.Sub(*lastCheckin)
	// если зашёл на следующий день
	if diff.Hours() >= 24 && diff.Hours() < 48 {
		// создаём sql запрос
		query, args, err := sq.Update("users").
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{"id": id}).
			Set("days_streak", streak+1).
			Set("last_checkin", now).
			ToSql()
		if err != nil {
			return false, err
		}
		// увеличиваем стрик и обновляем время захода
		_, err = pool.Exec(context.Background(), query, args...)
		if err != nil {
			return false, err
		}
		return true, err
	}
	// если игрок не заходил более двух суток
	if diff.Hours() >= 48 {
		// создаём sql запрос
		query, args, err := sq.Update("users").
			PlaceholderFormat(sq.Dollar).
			Set("days_streak", 1).
			Set("last_checkin", now).
			Where(sq.Eq{"id": id}).
			ToSql()
		if err != nil {
			return false, err
		}
		// обнуляем стрик
		_, err = pool.Exec(context.Background(), query, args...)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// если ни одно из условий не выполняется - игрок зашёл в течение суток; награды нет
	return false, nil
}

// Close закрывает соединение с бд
func Close() {
	pool.Close()
}
