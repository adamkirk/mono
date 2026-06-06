package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSSLMode string

const (
	PostgresSSLModeDisable    PostgresSSLMode = "disable"
	PostgresSSLModeAllow      PostgresSSLMode = "allow"
	PostgresSSLModePrefer     PostgresSSLMode = "prefer"
	PostgresSSLModeRequire    PostgresSSLMode = "require"
	PostgresSSLModeVerifyCA   PostgresSSLMode = "verify-ca"
	PostgresSSLModeVerifyFull PostgresSSLMode = "verify-full"
)

var sslModes []PostgresSSLMode = []PostgresSSLMode{
	PostgresSSLModeDisable,
	PostgresSSLModeAllow,
	PostgresSSLModePrefer,
	PostgresSSLModeRequire,
	PostgresSSLModeVerifyCA,
	PostgresSSLModeVerifyFull,
}

type PoolConfig struct {
	Host           string
	Username       string
	Password       string
	Port           uint16
	DBName         string          `yaml:"dbName"`
	MaxConnections int32           `yaml:"maxConnections"`
	MinConnections int32           `yaml:"minConnections"`
	SSLMode        PostgresSSLMode `yaml:"sslMode"`
}

func ISValidSSLMode(mode PostgresSSLMode) bool {
	for _, m := range sslModes {
		if mode == m {
			return true
		}
	}

	return false
}

func NewPool(poolConfig PoolConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		poolConfig.Username,
		poolConfig.Password,
		poolConfig.Host,
		poolConfig.Port,
		poolConfig.DBName,
		string(poolConfig.SSLMode),
	)

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(poolConfig.MaxConnections)
	config.MinConns = int32(poolConfig.MinConnections)

	return pgxpool.NewWithConfig(context.Background(), config)
}
