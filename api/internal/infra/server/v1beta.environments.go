package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/adamkirk/panoptes/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type V1BetaEnvironmentsController struct {
	environmentsHandler environmentsHandler
}

func (c *V1BetaEnvironmentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments",
		Summary:       "Create a new environment",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environments",
		},
		// Metadata: map[string]any{
		// 	OptDisableAllDefaultResponses:   false,
		// 	OptDisableDefaultAuthentication: false,
		// },
	}, ErrorHandler(c.Create, http.MethodPost))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.get", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments/{name}",
		Summary:       "Get an environment",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environments",
		},
		// Metadata: map[string]any{
		// 	OptDisableAllDefaultResponses:   false,
		// 	OptDisableDefaultAuthentication: false,
		// },
	}, ErrorHandler(c.Get, http.MethodGet))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.delete", string(v)),
		Method:        http.MethodDelete,
		Path:          "/environments/{name}",
		Summary:       "Delete an environment",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Environments",
		},
	}, ErrorHandler(c.Delete, http.MethodDelete))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.update", string(v)),
		Method:        http.MethodPatch,
		Path:          "/environments/{name}",
		Summary:       "Update an environment",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environments",
		},
	}, ErrorHandler(c.Update, http.MethodPatch))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.list", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments",
		Summary:       "List environments",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environments",
		},
	}, ErrorHandler(c.List, http.MethodGet))
}

func NewV1BetaEnvironmentsController(environmentsHandler environmentsHandler) *V1BetaEnvironmentsController {
	return &V1BetaEnvironmentsController{
		environmentsHandler: environmentsHandler,
	}
}

type V1BetaListEnvironmentsRequest struct {
	Page    int `query:"page" minimum:"1" default:"1" doc:"Page number, starting at 1."`
	PerPage int `query:"per_page" minimum:"1" maximum:"100" default:"20" doc:"Number of environments per page."`
}

func (req *V1BetaListEnvironmentsRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Page":
		return "query.page"
	case "PerPage":
		return "query.per_page"
	default:
		return targetField
	}
}

type V1BetaPaginationMeta struct {
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

type V1BetaListEnvironmentsMeta struct {
	Pagination V1BetaPaginationMeta `json:"pagination"`
}

type V1BetaListEnvironmentsResponseBody struct {
	Meta V1BetaListEnvironmentsMeta `json:"meta"`
	Data []V1BetaEnvironment        `json:"data"`
}

type V1BetaListEnvironmentsResponse struct {
	Status int
	Body   V1BetaListEnvironmentsResponseBody
}

func (c *V1BetaEnvironmentsController) List(ctx context.Context, req *V1BetaListEnvironmentsRequest) (*V1BetaListEnvironmentsResponse, error) {
	result, err := c.environmentsHandler.List(app.ListEnvironmentsDTO{
		Page:    req.Page,
		PerPage: req.PerPage,
	})

	if err != nil {
		return nil, err
	}

	data := make([]V1BetaEnvironment, len(result.Environments))
	for i, env := range result.Environments {
		data[i] = V1BetaEnvironment{
			ID:   env.ID.String(),
			Name: env.Name,
		}
	}

	return &V1BetaListEnvironmentsResponse{
		Status: http.StatusOK,
		Body: V1BetaListEnvironmentsResponseBody{
			Meta: V1BetaListEnvironmentsMeta{
				Pagination: V1BetaPaginationMeta{
					Total:      result.Total,
					TotalPages: result.TotalPages,
					Page:       result.Page,
					PerPage:    result.PerPage,
				},
			},
			Data: data,
		},
	}, nil
}

type V1BetaCreateEnvironmentRequestBody struct {
	Name string `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the environment. Must contain only alphanumeric characters and hyphens."`
}

type V1BetaCreateEnvironmentRequest struct {
	Body V1BetaCreateEnvironmentRequestBody
}

func (req *V1BetaCreateEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "body.name"
	default:
		return targetField
	}
}

type V1BetaEnvironment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type V1BetaEnvironmentResponseBody struct {
	Data V1BetaEnvironment `json:"data"`
}

type V1BetaEnvironmentResponse struct {
	Status int
	Body   V1BetaEnvironmentResponseBody
}

func (c *V1BetaEnvironmentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentRequest) (*V1BetaEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Create(app.CreateEnvironmentDTO{
		Name: req.Body.Name,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentResponse{
		Status: http.StatusCreated,
		Body: V1BetaEnvironmentResponseBody{
			Data: V1BetaEnvironment{
				ID:   env.ID.String(),
				Name: env.Name,
			},
		},
	}, nil
}

type V1BetaDeleteEnvironmentRequest struct {
	Name string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment to delete."`
}

func (req *V1BetaDeleteEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "path.name"
	default:
		return targetField
	}
}

type V1BetaDeleteEnvironmentResponse struct {
	Status int
}

func (c *V1BetaEnvironmentsController) Delete(ctx context.Context, req *V1BetaDeleteEnvironmentRequest) (*V1BetaDeleteEnvironmentResponse, error) {
	err := c.environmentsHandler.Delete(app.DeleteEnvironmentDTO{
		Name: req.Name,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaDeleteEnvironmentResponse{Status: http.StatusNoContent}, nil
}

type V1BetaUpdateEnvironmentRequestBody struct {
	Name string `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"New name for the environment. Must contain only alphanumeric characters and hyphens."`
}

type V1BetaUpdateEnvironmentRequest struct {
	Name string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Current name of the environment."`
	Body V1BetaUpdateEnvironmentRequestBody
}

func (req *V1BetaUpdateEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "CurrentName":
		return "path.name"
	case "Name":
		return "body.name"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentsController) Update(ctx context.Context, req *V1BetaUpdateEnvironmentRequest) (*V1BetaEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Update(app.UpdateEnvironmentDTO{
		CurrentName: req.Name,
		Name:        req.Body.Name,
	})

	if err != nil {
		return nil, err
	}

	if env == nil {
		// This shouldn't ever happen
		return nil, errors.New("no environment came back from update")
	}

	return &V1BetaEnvironmentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentResponseBody{
			Data: V1BetaEnvironment{
				ID:   env.ID.String(),
				Name: env.Name,
			},
		},
	}, nil
}

type V1BetaGetEnvironmentRequest struct {
	Name string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the environment. Must contain only alphanumeric characters and hyphens."`
}

func (req *V1BetaGetEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "path.name"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentsController) Get(ctx context.Context, req *V1BetaGetEnvironmentRequest) (*V1BetaEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Get(app.GetEnvironmentDTO{
		Name: req.Name,
	})

	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, huma.Error404NotFound("no environment with this name exists")
	}

	return &V1BetaEnvironmentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentResponseBody{
			Data: V1BetaEnvironment{
				ID:   env.ID.String(),
				Name: env.Name,
			},
		},
	}, nil
}
