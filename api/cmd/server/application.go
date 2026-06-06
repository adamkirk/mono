package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/adamkirk/panoptes/api/internal/app"
	"github.com/adamkirk/panoptes/api/internal/config"
	"github.com/adamkirk/panoptes/api/internal/infra/repository/postgres"
	"github.com/adamkirk/panoptes/api/internal/infra/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Application struct {
	// These should always be set, and safe to rely upon during construction of
	// any other services.
	stdout io.Writer
	stderr io.Writer
	logger *slog.Logger
	cfg    *config.Config

	pgPool                          *pgxpool.Pool
	environmentsRepository          *postgres.EnvironmentsRepository
	environmentsHandler             *app.EnvironmentsHandler
	environmentComponentsRepository *postgres.EnvironmentComponentsRepository
	environmentComponentsHandler    *app.EnvironmentComponentsHandler
	deploymentsRepository           *postgres.DeploymentsRepository
	deploymentsHandler              *app.DeploymentsHandler
	riverClient                     *river.Client[pgx.Tx]
}

func bindEnvs(v *viper.Viper, prefix string, t reflect.Type) {
	for field := range t.Fields() {
		key := field.Tag.Get("mapstructure")
		if key == "" {
			key = strings.ToLower(field.Name)
		}
		if prefix != "" {
			key = prefix + "." + key
		}
		if field.Type.Kind() == reflect.Struct {
			bindEnvs(v, key, field.Type)
		} else {
			_ = v.BindEnv(key)
		}
	}
}

func (a *Application) loadConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigFile("panoptes.server.yml")
	v.SetEnvPrefix("PANOPTES_SERVER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvs(v, "", reflect.TypeFor[config.Config]())

	_ = v.BindPFlag("logging.level", cmd.Root().PersistentFlags().Lookup("log-level"))
	_ = v.BindPFlag("logging.format", cmd.Root().PersistentFlags().Lookup("log-format"))

	if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	cfg := config.Default()
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	a.cfg = cfg

	return nil
}

func (a *Application) setupLogger() error {
	var level slog.Level
	if err := level.UnmarshalText([]byte(a.cfg.Logging.Level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", a.cfg.Logging.Level, err)
	}

	opts := &slog.HandlerOptions{Level: level}

	var l *slog.Logger

	switch a.cfg.Logging.Format {
	case "json":
		l = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	case "text":
		l = slog.New(slog.NewTextHandler(os.Stderr, opts))
	default:
		return fmt.Errorf("invalid log format %q: expected json or text", a.cfg.Logging.Format)
	}

	a.logger = l

	return nil
}

func (a *Application) GetServer() *server.Server {
	var accessLogger *slog.Logger

	if a.cfg.Server.AccessLogsEnabled {
		h := slog.NewJSONHandler(a.stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		accessLogger = slog.New(h)
		accessLogger = accessLogger.With("component", "http-server-access")
	}

	return server.New(
		a.cfg.GetServerPort(),
		a.logger.With("component", "http-server"),
		server.WithAccessLogger(accessLogger),
		server.WithApiVersionGroup(
			server.ApiVersionGroup{
				Version:     server.ApiVersionV1Beta,
				Controllers: a.GetV1BetaControllers(),
			},
		),
	)
}

func (a *Application) GetV1BetaControllers() []server.Controller {
	return []server.Controller{
		server.NewProbesController(),
		server.NewV1BetaEnvironmentsController(
			a.GetEnvironmentsHandler(),
		),
		server.NewV1BetaEnvironmentComponentsController(
			a.GetEnvironmentComponentsHandler(),
		),
		server.NewV1BetaEnvironmentComponentDeploymentsController(
			a.GetDeploymentsHandler(),
		),
	}
}

func (a *Application) GetPostgresPool() *pgxpool.Pool {
	if a.pgPool != nil {
		return a.pgPool
	}

	sslMode := postgres.PostgresSSLMode(a.cfg.DB.Postgres.SSLMode)

	if !postgres.ISValidSSLMode(sslMode) {
		cobra.CheckErr(fmt.Errorf("invalid ssl mode for postgres: %s", string(sslMode)))
	}

	p, err := postgres.NewPool(
		postgres.PoolConfig{
			Host:           a.cfg.DB.Postgres.Host,
			Username:       a.cfg.DB.Postgres.Username,
			Password:       a.cfg.DB.Postgres.Password,
			Port:           a.cfg.DB.Postgres.Port,
			DBName:         a.cfg.DB.Postgres.DBName,
			MaxConnections: a.cfg.DB.Postgres.MaxConnections,
			MinConnections: a.cfg.DB.Postgres.MinConnections,
			SSLMode:        sslMode,
		},
	)

	cobra.CheckErr(err)

	a.pgPool = p

	return a.pgPool
}

func (a *Application) GetRiverClient() *river.Client[pgx.Tx] {
	if a.riverClient != nil {
		return a.riverClient
	}

	client, err := postgres.NewRiverClient(a.GetPostgresPool())
	cobra.CheckErr(err)

	a.riverClient = client

	return a.riverClient
}

func (a *Application) GetEnvironmentsRepository() *postgres.EnvironmentsRepository {
	if a.environmentsRepository != nil {
		return a.environmentsRepository
	}

	a.environmentsRepository = postgres.NewEnvironmentsRepository(a.logger, a.GetPostgresPool(), a.GetRiverClient())

	return a.environmentsRepository
}

func (a *Application) GetEnvironmentComponentsRepository() *postgres.EnvironmentComponentsRepository {
	if a.environmentComponentsRepository != nil {
		return a.environmentComponentsRepository
	}

	a.environmentComponentsRepository = postgres.NewEnvironmentComponentsRepository(a.logger, a.GetPostgresPool(), a.GetRiverClient())

	return a.environmentComponentsRepository
}

func (a *Application) GetEnvironmentComponentsHandler() *app.EnvironmentComponentsHandler {
	if a.environmentComponentsHandler != nil {
		return a.environmentComponentsHandler
	}

	a.environmentComponentsHandler = app.NewEnvironmentComponentsHandler(
		a.logger,
		a.GetEnvironmentsRepository(),
		a.GetEnvironmentComponentsRepository(),
	)

	return a.environmentComponentsHandler
}

func (a *Application) GetDeploymentsRepository() *postgres.DeploymentsRepository {
	if a.deploymentsRepository != nil {
		return a.deploymentsRepository
	}

	a.deploymentsRepository = postgres.NewDeploymentsRepository(a.logger, a.GetPostgresPool(), a.GetRiverClient())

	return a.deploymentsRepository
}

func (a *Application) GetDeploymentsHandler() *app.DeploymentsHandler {
	if a.deploymentsHandler != nil {
		return a.deploymentsHandler
	}

	a.deploymentsHandler = app.NewDeploymentsHandler(
		a.logger,
		a.GetEnvironmentsRepository(),
		a.GetEnvironmentComponentsRepository(),
		a.GetDeploymentsRepository(),
	)

	return a.deploymentsHandler
}

func (a *Application) GetEnvironmentsHandler() *app.EnvironmentsHandler {
	if a.environmentsHandler != nil {
		return a.environmentsHandler
	}

	a.environmentsHandler = app.NewEnvironmentsHandler(a.logger, a.GetEnvironmentsRepository())

	return a.environmentsHandler
}

func NewApplication(cmd *cobra.Command) (*Application, error) {
	app := &Application{
		stderr: cmd.OutOrStderr(),
		stdout: cmd.OutOrStdout(),
	}

	if err := app.loadConfig(cmd); err != nil {
		return nil, err
	}

	if err := app.setupLogger(); err != nil {
		return nil, err
	}

	return app, nil
}
