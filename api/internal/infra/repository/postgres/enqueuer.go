package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type txEnqueuer struct {
	ctx    context.Context
	tx     pgx.Tx
	client *river.Client[pgx.Tx]
}

func (e *txEnqueuer) Enqueue(_ context.Context, args any) error {
	jobArgs, ok := args.(river.JobArgs)
	if !ok {
		return fmt.Errorf("enqueue: args must implement river.JobArgs (Kind() string)")
	}

	_, err := e.client.InsertTx(e.ctx, e.tx, jobArgs, nil)
	return err
}

func NewRiverClient(pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	return river.NewClient(riverpgxv5.New(pool), &river.Config{})
}
