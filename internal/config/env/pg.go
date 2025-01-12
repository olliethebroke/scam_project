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

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/pg.go - env variable %s not found", dsnEnvName))
	}
	return &pgConfig{
		dsn: dsn,
	}, nil
}
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
