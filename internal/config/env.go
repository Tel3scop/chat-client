package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

// Provider предоставляет интерфейс для получения конфигурации.
type Provider interface {
	Config() *Config
}

// Config общий конфиг
type Config struct {
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	Environment string `env:"ENV" envDefault:"local"`
	Encrypt     Encrypt
	Log         Log
	AuthService AuthService
	ChatService ChatService
}

// Config возвращаем сам конфиг
func (c Config) Config() *Config {
	return &c
}

// Encrypt конфиг с секретами
type Encrypt struct {
	RefreshTokenExpirationInMinutes int `env:"REFRESH_TOKEN_EXPIRATION" envDefault:"10"`
	RefreshTokenExpiration          time.Duration
	AccessTokenExpirationInMinutes  int `env:"ACCESS_TOKEN_EXPIRATION" envDefault:"2"`
	AccessTokenExpiration           time.Duration
	AuthPrefix                      string `env:"AUTH_PREFIX" envDefault:"Bearer "`
}

// Log конфиг для логов
type Log struct {
	FileName   string `env:"LOG_FILENAME" envDefault:"logs/app.log"`
	Level      string `env:"LOG_LEVEL" envDefault:"info"`
	MaxSize    int    `env:"LOG_MAXSIZE" envDefault:"5"`
	MaxBackups int    `env:"LOG_MAXBACKUPS" envDefault:"3"`
	MaxAge     int    `env:"LOG_MAXAGE" envDefault:"10"`
	Compress   bool   `env:"LOG_COMPRESS" envDefault:"false"`
}

// AuthService конфиг подключения к сервису авторизации
type AuthService struct {
	Host     string `env:"AUTH_HOST"`
	Port     int64  `env:"AUTH_PORT"`
	Protocol string `env:"AUTH_PROTOCOL"`
}

// ChatService конфиг подключения к сервису чатов
type ChatService struct {
	Host     string `env:"CHAT_HOST"`
	Port     int64  `env:"CHAT_PORT"`
	Protocol string `env:"CHAT_PROTOCOL"`
}

// New создаем новый конфиг
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	cfg.Encrypt.AccessTokenExpiration = time.Duration(cfg.Encrypt.AccessTokenExpirationInMinutes) * time.Minute
	cfg.Encrypt.RefreshTokenExpiration = time.Duration(cfg.Encrypt.RefreshTokenExpirationInMinutes) * time.Minute

	return cfg, nil
}
