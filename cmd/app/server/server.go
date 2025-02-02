package main

import (
	"crypto_scam/internal/api/access"
	admin_api "crypto_scam/internal/api/handler/admin"
	"crypto_scam/internal/api/handler/user"
	"crypto_scam/internal/config"
	"crypto_scam/internal/config/env"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"crypto_scam/internal/repository/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"os"
)

// TODO: 1) переделать логику SelectUserRole(должна брать данные из админки)
// TODO: 2) сделать DI контейнер
// TODO: 3) подумать над созданием интерфейса для логгирования

const (
	// роли
	creator = 2
	admin   = 1
	user    = 0

	createUserPostfix       = "/user/create/{id}"
	createFriendshipPostfix = "/friendship/create"

	// администраторские запросы
	adminCreateTaskPostfix = "/task/create"
	adminDeleteTaskPostfix = "/task/delete"
	adminGetTasksPostfix   = "/tasks/get"
	adminDeleteUserPostfix = "/user/delete"

	// пользовательские запросы
	getUserPostfix        = "/user/get"
	updateUserPostfix     = "/user/update"
	getLeaderboardPostfix = "/leaderboard/get"
	getUserTasksPostfix   = "/user/tasks/get"
)

var configPath string

func init() {
	configPath = os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "local.env"
	}
	logger.Info("path to config:", configPath)
}

func main() {
	// загружаем переменные из указанного .env файла в переменные окружения процесса
	if err := config.Load(configPath); err != nil {
		logger.Fatal("server.go/main - failed to load variables from env file: ", err)
	}
	// создаём новый конфиг хттп, чтобы создать слушающий адрес из глобальных переменных
	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		logger.Fatal("server.go/main - failed to create httpConfig: ", err)
	}
	// создаём конфиг постгреса, чтоб дсн нахуевертить из глобальных переменных
	pgConfig, err := env.NewPGConfig()
	if err != nil {
		logger.Fatal("server.go/main - failed to create pgConfig: ", err)
	}

	// создаём маршрутизатор и добавляем хендлеры
	r := chi.NewRouter()

	// если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	r.Use(middleware.Recoverer)

	// администраторские запросы
	r.With(access.UserAuthMiddleware(admin)).Get(adminGetTasksPostfix, admin_api.GetTasksHandler)
	r.With(access.UserAuthMiddleware(admin)).Post(adminCreateTaskPostfix, admin_api.CreateTaskHandler)
	r.With(access.UserAuthMiddleware(admin)).Delete(adminDeleteTaskPostfix, admin_api.DeleteTaskHandler)
	r.With(access.UserAuthMiddleware(admin)).Delete(adminDeleteUserPostfix, admin_api.DeleteUserHandler)

	// пользовательские запросы
	r.With(access.UserAuthMiddleware(user)).Get(getUserPostfix, user_api.GetUserHandler)
	r.With(access.UserAuthMiddleware(user)).Get(getLeaderboardPostfix, user_api.GetLeaderboardHandler)
	r.With(access.UserAuthMiddleware(user)).Get(getUserTasksPostfix, user_api.GetUserTasksHandler)
	r.With(access.UserAuthMiddleware(user)).Patch(updateUserPostfix, user_api.UpdateUserHandler)

	// коннектимся к бд
	db := connectToDatabase(pgConfig.DSN())

	// закрываем коннект с бд после завершения функции main
	defer db.Close()

	// запускаем сервер
	logger.Info("server is listening")
	err = http.ListenAndServe(httpConfig.Address(), r)
	if err != nil {
		logger.Fatal("server.go/main - failed while listening", err)
	}
}

// connectToDatabase устанавливает соединение с базой данных.
//
// Входным параметром функции является строка,
// содержащая dsn базы данных.
//
// Выходным параметром функции является переменная,
// реализующая интерфейс Repository.
func connectToDatabase(dsn string) repository.Repository {
	// инициализируем интерфейс взаимодействия с бд
	// в данном случае с postgresql
	var db repository.Repository = postgres.NewPostgres()

	// коннектимся к бд
	err := db.Connect(dsn)
	if err != nil {
		logger.Fatal("server.go/connectToDatabase - failed to connect to database: ", err)
	}

	// передаём интерфейс бд в апи
	admin_api.InitDatabase(db)
	user_api.InitDatabase(db)
	access.InitDatabase(db)

	// возвращаем реализацию интерфейса
	return db
}
