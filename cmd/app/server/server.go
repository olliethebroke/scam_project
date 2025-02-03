package main

import (
	"context"
	"crypto_scam/internal/app"
	"crypto_scam/internal/logger"
)

// TODO: 1) покрыть код тестами
// TODO: 2) настроить таймауты
// TODO: 3) настроить рейтлимитер

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
