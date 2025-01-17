package env

import (
	"crypto_scam/internal/config"
	"errors"
	"fmt"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const dsnEnvName = "PG_DSN"

type pgConfig struct {
	dsn string
}

// NewPGConfig создаёт новую конфигурацию postgres в зависимости от переменных окружения
func NewPGConfig() (*pgConfig, error) {
	// получаем данные из переменных окружения
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/pg.go - env variable %s not found", dsnEnvName))
	}
	// возвращаем конфиг
	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN возвращает ДСН базы данных
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
