package postgres

import (
	"context"
	"log/slog"

	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/adamkirk/panoptes/api/internal/infra/repository/postgres/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type EnvironmentsRepository struct {
	pool        *pgxpool.Pool
	l           *slog.Logger
	riverClient *river.Client[pgx.Tx]
}

func (r *EnvironmentsRepository) ByName(name string) (*common.Environment, error) {
	conn := db.New(r.pool)

	env, err := conn.GetEnvironmentByName(context.Background(), name)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		r.l.Error("failed to get environment", "error", err)
		return nil, err
	}

	return &common.Environment{
		ID:   env.ID.Bytes,
		Name: env.Name,
	}, nil
}

func (r *EnvironmentsRepository) Count() (int, error) {
	conn := db.New(r.pool)

	count, err := conn.CountEnvironments(context.Background())
	if err != nil {
		r.l.Error("failed to count environments", "error", err)
		return 0, err
	}

	return int(count), nil
}

func (r *EnvironmentsRepository) List(limit, offset int) ([]*common.Environment, error) {
	conn := db.New(r.pool)

	rows, err := conn.ListEnvironments(context.Background(), db.ListEnvironmentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		r.l.Error("failed to list environments", "error", err)
		return nil, err
	}

	envs := make([]*common.Environment, len(rows))
	for i, row := range rows {
		envs[i] = &common.Environment{
			ID:   row.ID.Bytes,
			Name: row.Name,
		}
	}

	return envs, nil
}

func (r *EnvironmentsRepository) Save(env *common.Environment, opts ...common.QueueJobOption) error {
	ctx := context.Background()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.l.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	_, err = db.New(tx).UpsertEnvironment(ctx, db.UpsertEnvironmentParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(env.ID[:]),
			Valid: true,
		},
		Name: env.Name,
	})
	if err != nil {
		r.l.Error("failed to save environment", "error", err)
		return err
	}

	enqueuer := &txEnqueuer{ctx: ctx, tx: tx, client: r.riverClient}
	for _, opt := range opts {
		if err := opt(enqueuer); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *EnvironmentsRepository) Delete(env *common.Environment) error {
	conn := db.New(r.pool)

	err := conn.DeleteEnvironmentByID(context.Background(), pgtype.UUID{
		Bytes: [16]byte(env.ID[:]),
		Valid: true,
	})
	if err != nil {
		r.l.Error("failed to delete environment", "error", err)
	}

	return err
}

func NewEnvironmentsRepository(l *slog.Logger, pool *pgxpool.Pool, riverClient *river.Client[pgx.Tx]) *EnvironmentsRepository {
	return &EnvironmentsRepository{
		pool:        pool,
		l:           l.With("component", "infra.postgres.environments_repository"),
		riverClient: riverClient,
	}
}
