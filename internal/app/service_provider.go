package app

import (
	"context"
	"crypto_scam/internal/closer"
	"crypto_scam/internal/config"
	"crypto_scam/internal/config/env"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"crypto_scam/internal/repository/postgres"
	"crypto_scam/pkg/hooks/telegram"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	tgConfig   config.TGConfig
	httpConfig config.HTTPConfig

	db repository.Repository
}

// newServiceProvider создаёт новый serviceProvider.
//
// Возвращает указатель на созданный объект типа serviceProvider.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}

}

// PGConfig возвращает переменную, реализующую интерфейс PGConfig
// для конфигурации работы с базой данных PostgreSQL. Если такая переменная
// не была инициализирована, то этот метод её инциализирует.
//
// Выходным параметром метода является объект, реализующий интерфейс
// PGConfig
func (s *serviceProvider) PGConfig() config.PGConfig {
	// если интерфейс не реализован
	if s.pgConfig == nil {
		// создаём новый конфиг
		// реализуем интерфейс
		cfg, err := env.NewPGConfig()
		if err != nil {
			logger.Fatal("failed to create pgConfig")
		}

		// передаём переменную, реализующую интерфейс,
		// в поле экземляра структуры serviceProvider
		s.pgConfig = cfg
	}

	// возвращаем переменную, реализующую интерфейс
	return s.pgConfig
}

// TGConfig возвращает переменную, реализующую интерфейс TGConfig
// для конфигурации работы с Telegram-API. Если такая переменная
// не была инициализирована, то этот метод её инциализирует.
//
// Выходным параметром метода является объект, реализующий интерфейс
// TGConfig
func (s *serviceProvider) TGConfig() config.TGConfig {
	// если интерфейс не реализован
	if s.tgConfig == nil {
		// создаём новый конфиг
		// реализуем интерфейс
		cfg, err := env.NewTgConfig()
		if err != nil {
			logger.Fatal("failed to create tgConfig: ", err)
		}

		// передаём переменную, реализующую интерфейс,
		// в поле экземляра структуры serviceProvider
		// и инициализируем переменную конфигурации в пакете telegram
		s.tgConfig = cfg
		telegram.InitTGConfig(cfg)
	}

	// возвращаем переменную, реализующую интерфейс
	return s.tgConfig
}

// HTTPConfig возвращает переменную, реализующую интерфейс HTTPConfig
// для конфигурации работы с протоколом HTTP. Если такая переменная
// не была инициализирована, то этот метод её инциализирует.
//
// Выходным параметром метода является объект, реализующий интерфейс
// HTTPConfig
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	// если интерфейс не реализован
	if s.httpConfig == nil {
		// создаём новый конфиг
		// реализуем интерфейс
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			logger.Fatal("failed to create httpConfig")
		}

		// передаём переменную, реализующую интерфейс,
		// в поле экземляра структуры serviceProvider
		s.httpConfig = cfg
	}

	// возвращаем переменную, реализующую интерфейс
	return s.httpConfig
}

// DB возвращает переменную, реализующую интерфейс Repository
// для работы с базой данных. Если такая переменная
// не была инициализирована, то этот метод её инциализирует.
//
// Выходным параметром метода является объект, реализующий интерфейс
// Repository
func (s *serviceProvider) DB(ctx context.Context) repository.Repository {
	// если интерфейс не реализован
	if s.db == nil {
		// создаём новый объект, реализующий интерфейс
		// взаимодействия с базой данных.
		// подключаемся к бд
		db, err := postgres.NewPostgres(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create database: ", err)
		}

		// закрываем пул соединений с базой данных
		// после завершения работы сервера
		closer.Add(func() error {
			db.Close()
			return nil
		})

		// передаём переменную, реализующую интерфейс,
		// в поле экземляра структуры serviceProvider
		s.db = db
	}

	// возвращаем переменную, реализующую интерфейс
	return s.db
}
