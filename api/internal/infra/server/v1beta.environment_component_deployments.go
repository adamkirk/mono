package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adamkirk/panoptes/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type V1BetaEnvironmentComponentDeploymentsController struct {
	deploymentsHandler deploymentsHandler
}

func (c *V1BetaEnvironmentComponentDeploymentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.deployments.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments/{environment_name}/components/{component_name}/deployments",
		Summary:       "Create a new deployment for an environment component",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environment Component Deployments",
		},
	}, ErrorHandler(c.Create, http.MethodPost))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.deployments.get", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments/{environment_name}/components/{component_name}/deployments/{deployment_id}",
		Summary:       "Get a deployment by ID",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environment Component Deployments",
		},
	}, ErrorHandler(c.Get, http.MethodGet))
}

func NewV1BetaEnvironmentComponentDeploymentsController(deploymentsHandler deploymentsHandler) *V1BetaEnvironmentComponentDeploymentsController {
	return &V1BetaEnvironmentComponentDeploymentsController{
		deploymentsHandler: deploymentsHandler,
	}
}

// shared response types

type V1BetaEnvironmentComponentDeployment struct {
	ID                       string    `json:"id"`
	Status                   string    `json:"status"`
	CreatedAt                time.Time `json:"created_at"`
	EnvironmentName          string    `json:"environment_name"`
	EnvironmentComponentName string    `json:"environment_component_name"`
}

type V1BetaEnvironmentComponentDeploymentResponseBody struct {
	Data V1BetaEnvironmentComponentDeployment `json:"data"`
}

type V1BetaEnvironmentComponentDeploymentResponse struct {
	Status int
	Body   V1BetaEnvironmentComponentDeploymentResponseBody
}

// create

type V1BetaCreateEnvironmentComponentDeploymentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	ComponentName   string `path:"component_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the component."`
}

func (req *V1BetaCreateEnvironmentComponentDeploymentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "ComponentName":
		return "path.component_name"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentComponentDeploymentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentComponentDeploymentRequest) (*V1BetaEnvironmentComponentDeploymentResponse, error) {
	deployment, err := c.deploymentsHandler.Create(app.CreateDeploymentDTO{
		EnvironmentName: req.EnvironmentName,
		ComponentName:   req.ComponentName,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentComponentDeploymentResponse{
		Status: http.StatusCreated,
		Body: V1BetaEnvironmentComponentDeploymentResponseBody{
			Data: V1BetaEnvironmentComponentDeployment{
				ID:                       deployment.ID.String(),
				Status:                   string(deployment.Status),
				CreatedAt:                deployment.CreatedAt,
				EnvironmentName:          req.EnvironmentName,
				EnvironmentComponentName: req.ComponentName,
			},
		},
	}, nil
}

// get

type V1BetaGetEnvironmentComponentDeploymentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	ComponentName   string `path:"component_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the component."`
	DeploymentID    string `path:"deployment_id" minLength:"1" doc:"ID of the deployment."`
}

func (req *V1BetaGetEnvironmentComponentDeploymentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "ComponentName":
		return "path.component_name"
	case "DeploymentID":
		return "path.deployment_id"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentComponentDeploymentsController) Get(ctx context.Context, req *V1BetaGetEnvironmentComponentDeploymentRequest) (*V1BetaEnvironmentComponentDeploymentResponse, error) {
	id, err := uuid.Parse(req.DeploymentID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("invalid deployment ID", &huma.ErrorDetail{
			Message:  "must be a valid UUID",
			Location: "path.deployment_id",
			Value:    req.DeploymentID,
		})
	}

	deployment, err := c.deploymentsHandler.Get(app.GetDeploymentDTO{
		EnvironmentName: req.EnvironmentName,
		ComponentName:   req.ComponentName,
		DeploymentID:    id,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentComponentDeploymentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentComponentDeploymentResponseBody{
			Data: V1BetaEnvironmentComponentDeployment{
				ID:                       deployment.ID.String(),
				Status:                   string(deployment.Status),
				CreatedAt:                deployment.CreatedAt,
				EnvironmentName:          req.EnvironmentName,
				EnvironmentComponentName: req.ComponentName,
			},
		},
	}, nil
}
