package config

type LoggingConfig struct {
	Level  string
	Format string
}

type ServerConfig struct {
	Port              int
	AccessLogsEnabled bool
}

type PostgresConfig struct {
	Host           string
	Username       string
	Password       string
	Port           uint16
	DBName         string `mapstructure:"dbName"`
	MaxConnections int32  `mapstructure:"maxConnections"`
	MinConnections int32  `mapstructure:"minConnections"`
	SSLMode        string `mapstructure:"sslMode"`
}
type DBConfig struct {
	Postgres PostgresConfig
}
type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
	DB      DBConfig
}

func (c *Config) GetServerPort() int {
	return c.Server.Port
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:              8080,
			AccessLogsEnabled: true,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
