package main

import (
	"context"
	"crypto_scam/internal/app"
	"crypto_scam/internal/logger"
)

func main() {
	// инициализируем зависимости для работы приложения
	a, err := app.NewApp(context.Background())
	if err != nil {
		logger.Fatal("failed to initialize application: ", err)
	}

	// запускаем приложение
	err = a.Run()
	if err != nil {
		logger.Fatal("failed to run application: ", err)
	}
}
