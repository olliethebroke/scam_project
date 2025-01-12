package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

// HTTPConfig интерфейс конфигуратора слушающего интерфейса
type HTTPConfig interface {
	Address() string
}

// PGConfig интерфейс конфигуратора базы данных
type PGConfig interface {
	DSN() string
}

type TGConfig interface {
	Token() string
	ChatId() int64
}

// Load загружает переменные из env файла в переменные окружения процесса
func Load(path string) error {
	err := godotenv.Load(fmt.Sprintf("/%s", path))
	if err != nil {
		return err
	}
	return nil
}
