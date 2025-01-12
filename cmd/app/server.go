package main

import (
	"crypto_scam/internal/api"
	"crypto_scam/internal/config"
	"crypto_scam/internal/config/env"
	"crypto_scam/internal/logger/logrus_logger"
	"crypto_scam/internal/storage/postgres"
	"github.com/go-chi/chi"
	"net/http"
	"os"
)

const (
	createUserPostfix     = "/user/create/{id}"
	getUserPostfix        = "/user/get/{id}"
	updateUserPostfix     = "/user/update/{id}"
	getLeaderboardPostfix = "/leaderboard/get"
)

var configPath string

func init() {
	configPath = os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "local.env"
	}
	logrus_logger.Log.Info("path to config:", configPath)
}

func main() {
	// загружаем переменные из указанного .env файла в переменные окружения процесса
	if err := config.Load(configPath); err != nil {
		logrus_logger.Log.Fatal("cmd/app/server.go - failed to load variables from env file: ", err)
	}
	// создаём новый конфиг хттп, чтобы создать слушающий адрес из глобальных переменных
	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		logrus_logger.Log.Fatal("cmd/app/server.go - failed to create httpConfig: ", err)
	}
	// создаём конфиг постгреса, чтоб дсн нахуевертить из глобальных переменных
	pgConfig, err := env.NewPGConfig()
	if err != nil {
		logrus_logger.Log.Fatal("cmd/app/server.go - failed to create pgConfig: ", err)
	}
	// создаём маршрутизатор и добавляем ручки
	r := chi.NewRouter()
	r.Post(createUserPostfix, api.CreateUserHandler)
	r.Get(getUserPostfix, api.GetUserHandler)
	r.Get(getLeaderboardPostfix, api.GetLeaderboardHandler)
	r.Put(updateUserPostfix, api.UpdateUserHandler)

	// коннектимся к бд
	err = postgres.Connect(pgConfig.DSN())
	if err != nil {
		logrus_logger.Log.Fatal("cmd/app/server.go - failed to connect to database: ", err)
	}
	defer postgres.Close()
	logrus_logger.Log.Info("server is listening")
	// запускаем сервер
	err = http.ListenAndServe(httpConfig.Address(), r)
	if err != nil {
		logrus_logger.Log.Fatal("cmd/app/server.go - failed while listening", err)
	}
}
