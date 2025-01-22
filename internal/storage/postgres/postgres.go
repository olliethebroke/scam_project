package postgres

import (
	"context"
	"crypto_scam/internal/storage/models"
	"crypto_scam/pkg/hooks/telegram"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"slices"
	"time"
)

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
func InsertUser(id int64, firstname string) (*models.SelectUserResponse, error) {
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
		Columns("id, firstname, is_premium").
		Values(id, firstname, isPremium).
		Suffix("RETURNING firstname").
		ToSql()
	if err != nil {
		return nil, err
	}
	// вносим пользователя в бд
	err = pool.QueryRow(context.Background(), query, args...).Scan(&firstname)
	if err != nil {
		return nil, err
	}
	// возвращаем данные нового пользователя
	return &models.SelectUserResponse{
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

// SelectUser считывает данные о пользователе из бд
// возвращает ошибку и данные пользователя
func SelectUser(id int64) (*models.SelectUserResponse, error) {
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
	var info = &models.SelectUserResponse{}
	var lastCheckin = &time.Time{}
	// получаем данные о пользователе
	err = pool.QueryRow(context.Background(), query, args...).Scan(
		&info.Id, &info.Firstname,
		&info.Blocks, &info.Record, lastCheckin,
		&info.DaysStreak, &info.InvitedFriends,
		&info.IsPremium, &info.League)
	if err != nil {
		return nil, err
	}
	// обновляем серию ежедневных заходов
	award, err := updateStreak(id, lastCheckin, &info.DaysStreak)
	if err != nil {
		return nil, err
	}
	// присваиваем переменной award в респонсе, нужно ли выдавать игроку награду за вход
	info.Award = award
	return info, nil
}

// updateStreak обновляет стрик ежедневных заходов в игру;
// аргумент lastCheckin - время последнего захода, за который давалась daily награда;
// аргумент streak - текущий стрик игрока с указанным id;
// возвращает ошибку и флаг о том, нужна ли награда
func updateStreak(id int64, lastCheckin *time.Time, streak *int) (bool, error) {
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
		_, err = pool.Exec(context.Background(), query, args...)
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
		_, err = pool.Exec(context.Background(), query, args...)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// если ни одно из условий не выполняется - игрок зашёл в течение суток; награды нет
	return false, nil
}

// SelectLeaders считывает данные о лидерах из бд;
// возвращает ошибку и мапу, содержащую для каждой лиги(int16) слайс из лидеров
func SelectLeaders() (map[int16][]*models.Leader, error) {
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
	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	// создаём мапу лидеров в лигах
	leaders := make(map[int16][]*models.Leader)
	// проходимся по каждой строке, чтобы считать данные
	for rows.Next() {
		leader := &models.Leader{}
		var league int16
		// считываем данные о лидере
		err = rows.Scan(&leader.Id, &leader.Position, &leader.Firstname, &leader.Blocks, &league)
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
func UpdateUser(id int64, info *models.UpdateUserRequest) error {
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

// InsertFriendship создаёт новую дружбу в бд и увеличивает счётчик друзей у приглашающего пользователя;
// принимает на вход два id - приглашённого и приглашающего пользователей
func InsertFriendship(invitedUserId int64, invitingUserId int64) error {
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
	_, err = pool.Exec(context.Background(), query, args...)
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
	_, err = pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return nil
}

// SelectUserTasks возвращает из бд массив заданий пользователя по его id
func SelectUserTasks(id int64) ([]*models.GetUserTaskResponse, error) {
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
	rows, err := pool.Query(context.Background(), query, args...)
	// создаём слайс выполненных заданий
	completedTasks := make([]int, 0, 6)
	if err != nil {
		return nil, err
	}
	// отмечаем выполненные пользователем задания в слайсе completedTasks
	for rows.Next() {
		var completedTaskId int
		// считываем id выполненного задания
		err = rows.Scan(&completedTaskId)
		if err != nil {
			return nil, err
		}
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
	rows, err = pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	// создаём слайс заданий
	tasks := make([]*models.GetUserTaskResponse, 0, 6)
	// проходимся по полученному списку и добавляем каждое задание в слайс
	for rows.Next() {
		t := &models.GetUserTaskResponse{}
		// считываем информацию о задании
		err = rows.Scan(&t.Task.Id, &t.Task.Description, &t.Task.Reward, &t.Task.ActionType, &t.Task.ActionData)
		if err != nil {
			return nil, err
		}
		// если задание есть в списке выполненных, то ставим флаг
		t.IsCompleted = slices.Contains(completedTasks, t.Task.Id)
		// добавляем задание в слайс
		tasks = append(tasks, t)
	}
	// возвращаем результат
	return tasks, nil
}

// InsertTask добавляет новое задание в бд
func InsertTask(task *models.Task) error {
	// создаём sql запрос
	query, args, err := sq.Insert("tasks").
		PlaceholderFormat(sq.Dollar).
		Columns("description", "reward", "action_type", "action_data").
		Values(task.Description, task.Reward, task.ActionType, task.ActionData).
		ToSql()
	if err != nil {
		return err
	}
	// добавляем задание в бд
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
