package repository

import (
	"context"
	"crypto_scam/internal/repository/model"
)

// Repository имплементирует взаимодействие с базой данных.
//
// Close - закрывает соединение с бд.
//
// InsertUser - добавляет пользователя,
// SelectUser - считывает пользователя,
// SelectUserRole - считывает роль пользователя,
// UpdateUser - обновляет информацию о пользователе,
// DeleteUser - удаляет пользователя.
//
// InsertLeader - добавляет лидера в лигу | test,
// SelectLeaders - считывает по 100 лидеров из каждой лиги.
//
// InsertTask - добавляет задание,
// InsertCompletedTask - добавляет выполненное пользователем задание,
// SelectTask - считывает задание,
// SelectTasks - считывает все задания,
// SelectUserTasks - считывает задания пользователя,
// DeleteTask - удаляет задание.
//
// InsertFriendship - добавляет дружбу.
type Repository interface {
	Close()

	InsertUser(ctx context.Context, id int64, firstname string) (*model.User, error)
	SelectUser(ctx context.Context, id int64) (*model.User, error)
	SelectUserRole(ctx context.Context, id int64) (*model.UserRole, error)
	UpdateUser(ctx context.Context, id int64, info *model.Update) error
	DeleteUser(ctx context.Context, id int64) error

	InsertLeader(ctx context.Context, leader *model.Leader, league int16) error
	SelectLeaders(ctx context.Context) (map[int16][]*model.Leader, error)

	InsertTask(ctx context.Context, task *model.Task) (*model.Task, error)
	InsertCompletedTask(ctx context.Context, userId int64, taskId int16) error
	SelectTask(ctx context.Context, id int16) (*model.Task, error)
	SelectTasks(ctx context.Context) ([]*model.Task, error)
	SelectUserTasks(ctx context.Context, id int64) ([]*model.Task, error)
	DeleteTask(ctx context.Context, id int16) error

	InsertFriendship(ctx context.Context, invitedUserId int64, invitingUserId int64) error

	InsertAdmin(ctx context.Context, id int64) error
}
