package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamkirk/panoptes/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type V1BetaEnvironmentComponentsController struct {
	environmentComponentsHandler environmentComponentsHandler
}

func (c *V1BetaEnvironmentComponentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments/{environment_name}/components",
		Summary:       "Create a new environment component",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.Create, http.MethodPost))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.delete", string(v)),
		Method:        http.MethodDelete,
		Path:          "/environments/{environment_name}/components/{name}",
		Summary:       "Delete an environment component",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.Delete, http.MethodDelete))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.update", string(v)),
		Method:        http.MethodPatch,
		Path:          "/environments/{environment_name}/components/{name}",
		Summary:       "Update an environment component",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.Update, http.MethodPatch))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.get", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments/{environment_name}/components/{name}",
		Summary:       "Get an environment component",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.Get, http.MethodGet))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.list", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments/{environment_name}/components",
		Summary:       "List environment components",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.List, http.MethodGet))
}

func NewV1BetaEnvironmentComponentsController(environmentComponentsHandler environmentComponentsHandler) *V1BetaEnvironmentComponentsController {
	return &V1BetaEnvironmentComponentsController{
		environmentComponentsHandler: environmentComponentsHandler,
	}
}

type V1BetaCreateEnvironmentComponentRequestBody struct {
	Name          string `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the component within the environment."`
	ChartName     string `json:"chart_name" minLength:"1" doc:"Name of the Helm chart."`
	ChartVersion  string `json:"chart_version" minLength:"1" doc:"Version of the Helm chart."`
	ChartRegistry string `json:"chart_registry" minLength:"1" doc:"Registry where the Helm chart is hosted."`
}

type V1BetaCreateEnvironmentComponentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Body            V1BetaCreateEnvironmentComponentRequestBody
}

func (req *V1BetaCreateEnvironmentComponentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "Name":
		return "body.name"
	case "ChartName":
		return "body.chart_name"
	case "ChartVersion":
		return "body.chart_version"
	case "ChartRegistry":
		return "body.chart_registry"
	default:
		return targetField
	}
}

type V1BetaDeleteEnvironmentComponentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Name            string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the component to delete."`
}

func (req *V1BetaDeleteEnvironmentComponentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "Name":
		return "path.name"
	default:
		return targetField
	}
}

type V1BetaDeleteEnvironmentComponentResponse struct {
	Status int
}

func (c *V1BetaEnvironmentComponentsController) Delete(ctx context.Context, req *V1BetaDeleteEnvironmentComponentRequest) (*V1BetaDeleteEnvironmentComponentResponse, error) {
	err := c.environmentComponentsHandler.Delete(app.DeleteEnvironmentComponentDTO{
		EnvironmentName: req.EnvironmentName,
		Name:            req.Name,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaDeleteEnvironmentComponentResponse{Status: http.StatusNoContent}, nil
}

type V1BetaUpdateEnvironmentComponentRequestBody struct {
	Name          *string `json:"name,omitempty" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"New name for the component."`
	ChartName     *string `json:"chart_name,omitempty" minLength:"1" doc:"Name of the Helm chart."`
	ChartVersion  *string `json:"chart_version,omitempty" minLength:"1" doc:"Version of the Helm chart."`
	ChartRegistry *string `json:"chart_registry,omitempty" minLength:"1" doc:"Registry where the Helm chart is hosted."`
}

type V1BetaUpdateEnvironmentComponentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Name            string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Current name of the component."`
	Body            V1BetaUpdateEnvironmentComponentRequestBody
}

func (req *V1BetaUpdateEnvironmentComponentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "CurrentName":
		return "path.name"
	case "Name":
		return "body.name"
	case "ChartName":
		return "body.chart_name"
	case "ChartVersion":
		return "body.chart_version"
	case "ChartRegistry":
		return "body.chart_registry"
	case "Body":
		return "body"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentComponentsController) Update(ctx context.Context, req *V1BetaUpdateEnvironmentComponentRequest) (*V1BetaEnvironmentComponentResponse, error) {
	component, err := c.environmentComponentsHandler.Update(app.UpdateEnvironmentComponentDTO{
		EnvironmentName: req.EnvironmentName,
		CurrentName:     req.Name,
		Name:            req.Body.Name,
		ChartName:       req.Body.ChartName,
		ChartVersion:    req.Body.ChartVersion,
		ChartRegistry:   req.Body.ChartRegistry,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentComponentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentComponentResponseBody{
			Data: V1BetaEnvironmentComponent{
				ID:            component.ID.String(),
				EnvironmentID: component.EnvironmentID.String(),
				Name:          component.Name,
				ChartName:     component.ChartName,
				ChartVersion:  component.ChartVersion,
				ChartRegistry: component.ChartRegistry,
			},
		},
	}, nil
}

type V1BetaGetEnvironmentComponentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Name            string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the component."`
}

func (req *V1BetaGetEnvironmentComponentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "Name":
		return "path.name"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentComponentsController) Get(ctx context.Context, req *V1BetaGetEnvironmentComponentRequest) (*V1BetaEnvironmentComponentResponse, error) {
	component, err := c.environmentComponentsHandler.Get(app.GetEnvironmentComponentDTO{
		EnvironmentName: req.EnvironmentName,
		Name:            req.Name,
	})

	if err != nil {
		return nil, err
	}

	if component == nil {
		return nil, huma.Error404NotFound("no component with this name exists in this environment")
	}

	return &V1BetaEnvironmentComponentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentComponentResponseBody{
			Data: V1BetaEnvironmentComponent{
				ID:            component.ID.String(),
				EnvironmentID: component.EnvironmentID.String(),
				Name:          component.Name,
				ChartName:     component.ChartName,
				ChartVersion:  component.ChartVersion,
				ChartRegistry: component.ChartRegistry,
			},
		},
	}, nil
}

type V1BetaListEnvironmentComponentsRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Page            int    `query:"page" minimum:"1" default:"1" doc:"Page number, starting at 1."`
	PerPage         int    `query:"per_page" minimum:"1" maximum:"100" default:"20" doc:"Number of components per page."`
}

func (req *V1BetaListEnvironmentComponentsRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "Page":
		return "query.page"
	case "PerPage":
		return "query.per_page"
	default:
		return targetField
	}
}

type V1BetaListEnvironmentComponentsResponseBody struct {
	Meta V1BetaListEnvironmentsMeta   `json:"meta"`
	Data []V1BetaEnvironmentComponent `json:"data"`
}

type V1BetaListEnvironmentComponentsResponse struct {
	Status int
	Body   V1BetaListEnvironmentComponentsResponseBody
}

func (c *V1BetaEnvironmentComponentsController) List(ctx context.Context, req *V1BetaListEnvironmentComponentsRequest) (*V1BetaListEnvironmentComponentsResponse, error) {
	result, err := c.environmentComponentsHandler.List(app.ListEnvironmentComponentsDTO{
		EnvironmentName: req.EnvironmentName,
		Page:            req.Page,
		PerPage:         req.PerPage,
	})

	if err != nil {
		return nil, err
	}

	data := make([]V1BetaEnvironmentComponent, len(result.Components))
	for i, c := range result.Components {
		data[i] = V1BetaEnvironmentComponent{
			ID:            c.ID.String(),
			EnvironmentID: c.EnvironmentID.String(),
			Name:          c.Name,
			ChartName:     c.ChartName,
			ChartVersion:  c.ChartVersion,
			ChartRegistry: c.ChartRegistry,
		}
	}

	return &V1BetaListEnvironmentComponentsResponse{
		Status: http.StatusOK,
		Body: V1BetaListEnvironmentComponentsResponseBody{
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

type V1BetaEnvironmentComponent struct {
	ID            string `json:"id"`
	EnvironmentID string `json:"environment_id"`
	Name          string `json:"name"`
	ChartName     string `json:"chart_name"`
	ChartVersion  string `json:"chart_version"`
	ChartRegistry string `json:"chart_registry"`
}

type V1BetaEnvironmentComponentResponseBody struct {
	Data V1BetaEnvironmentComponent `json:"data"`
}

type V1BetaEnvironmentComponentResponse struct {
	Status int
	Body   V1BetaEnvironmentComponentResponseBody
}

func (c *V1BetaEnvironmentComponentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentComponentRequest) (*V1BetaEnvironmentComponentResponse, error) {
	component, err := c.environmentComponentsHandler.Create(app.CreateEnvironmentComponentDTO{
		EnvironmentName: req.EnvironmentName,
		Name:            req.Body.Name,
		ChartName:       req.Body.ChartName,
		ChartVersion:    req.Body.ChartVersion,
		ChartRegistry:   req.Body.ChartRegistry,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentComponentResponse{
		Status: http.StatusCreated,
		Body: V1BetaEnvironmentComponentResponseBody{
			Data: V1BetaEnvironmentComponent{
				ID:            component.ID.String(),
				EnvironmentID: component.EnvironmentID.String(),
				Name:          component.Name,
				ChartName:     component.ChartName,
				ChartVersion:  component.ChartVersion,
				ChartRegistry: component.ChartRegistry,
			},
		},
	}, nil
}
