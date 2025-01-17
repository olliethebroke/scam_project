package env

import (
	"crypto_scam/internal/config"
	"errors"
	"fmt"
	"net"
	"os"
)

var _ config.HTTPConfig = (*httpConfig)(nil)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig создаёт новую конфигурацию http соединения в зависимости от переменных окружения
func NewHTTPConfig() (*httpConfig, error) {
	// получаем данные из переменных окружения
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/http.go - env variable %s not found", httpHostEnvName))
	}
	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/http.go - env variable %s not found", httpPortEnvName))
	}
	// возвращаем конфиг
	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

// Address возвращает интерфейс
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
