package app

import (
	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/google/uuid"
)

type environmentsRepository interface {
	ByName(name string) (*common.Environment, error)
	List(limit, offset int) ([]*common.Environment, error)
	Count() (int, error)
	Save(env *common.Environment, opts ...common.QueueJobOption) error
	Delete(env *common.Environment) error
}

type environmentComponentsRepository interface {
	ByEnvironmentAndName(environmentID uuid.UUID, name string) (*common.EnvironmentComponent, error)
	CountByEnvironment(environmentID uuid.UUID) (int, error)
	ListByEnvironment(environmentID uuid.UUID, limit, offset int) ([]*common.EnvironmentComponent, error)
	Save(c *common.EnvironmentComponent, opts ...common.QueueJobOption) error
	Delete(c *common.EnvironmentComponent) error
}

type deploymentsRepository interface {
	ByID(id uuid.UUID) (*common.Deployment, error)
	Save(d *common.Deployment, opts ...common.QueueJobOption) error
}
