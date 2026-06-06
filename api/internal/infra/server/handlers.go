package server

import (
	"github.com/adamkirk/panoptes/api/internal/app"
	"github.com/adamkirk/panoptes/api/internal/common"
)

type environmentComponentsHandler interface {
	Create(dto app.CreateEnvironmentComponentDTO) (*common.EnvironmentComponent, error)
	Get(dto app.GetEnvironmentComponentDTO) (*common.EnvironmentComponent, error)
	Update(dto app.UpdateEnvironmentComponentDTO) (*common.EnvironmentComponent, error)
	Delete(dto app.DeleteEnvironmentComponentDTO) error
	List(dto app.ListEnvironmentComponentsDTO) (*app.ListEnvironmentComponentsResult, error)
}

type environmentsHandler interface {
	Create(dto app.CreateEnvironmentDTO) (*common.Environment, error)
	Get(dto app.GetEnvironmentDTO) (*common.Environment, error)
	List(dto app.ListEnvironmentsDTO) (*app.ListEnvironmentsResult, error)
	Update(dto app.UpdateEnvironmentDTO) (*common.Environment, error)
	Delete(dto app.DeleteEnvironmentDTO) error
}

type deploymentsHandler interface {
	Create(dto app.CreateDeploymentDTO) (*common.Deployment, error)
	Get(dto app.GetDeploymentDTO) (*common.Deployment, error)
}
