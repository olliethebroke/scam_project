package repository

import "crypto_scam/internal/repository/model"

// Repository имплементирует взаимодействие с базой данных.
//
// Connect - создаёт пул соединений с бд,
// Close - закрывает соединение с бд.
//
// InsertUser - добавляет пользователя,
// SelectUser - считывает пользователя,
// SelectUserRole - считывает роль пользователя,
// UpdateUser - обновляет информацию о пользователе,
// DeleteUser - удаляет пользователя.
//
// SelectLeaders - считывает по 100 лидеров из каждой лиги.
//
// InsertTask - добавляет задание,
// SelectTasks - считывает все задания,
// SelectUserTasks - считывает задания пользователя,
// DeleteTask - удаляет задание.
//
// InsertFriendship - добавляет дружбу.
type Repository interface {
	Connect(dsn string) error
	Close()

	InsertUser(id int64, firstname string) (*model.User, error)
	SelectUser(id int64) (*model.User, error)
	SelectUserRole(id int64) (*model.UserRole, error)
	UpdateUser(id int64, info *model.Update) error
	DeleteUser(id int64) error

	SelectLeaders() (map[int16][]*model.Leader, error)

	InsertTask(task *model.Task) error
	SelectTasks() ([]*model.Task, error)
	SelectUserTasks(id int64) ([]*model.Task, error)
	DeleteTask(id int16) error

	InsertFriendship(invitedUserId int64, invitingUserId int64) error
}
