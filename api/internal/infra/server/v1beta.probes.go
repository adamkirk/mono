package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type V1BetaProbesStartupRequest struct{}

type V1BetaProbesController struct{}

func (c *V1BetaProbesController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.probes.startup", string(v)),
		Method:        http.MethodGet,
		Path:          "/_/probes/startup",
		Summary:       "Check if the app is started up",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Probes",
		},
		Metadata: map[string]any{
			OptDisableAllDefaultResponses:   true,
			OptDisableDefaultAuthentication: true,
		},
	}, ErrorHandler(c.Startup, http.MethodGet))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.probes.ready", string(v)),
		Method:        http.MethodGet,
		Path:          "/_/probes/ready",
		Summary:       "Check if the app is ready to serve connections",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Probes",
		},
		Metadata: map[string]any{
			OptDisableAllDefaultResponses:   true,
			OptDisableDefaultAuthentication: true,
		},
	}, ErrorHandler(c.Ready, http.MethodGet))
}

func NewProbesController() *V1BetaProbesController {
	return &V1BetaProbesController{}
}

func (c *V1BetaProbesController) Startup(ctx context.Context, req *V1BetaProbesStartupRequest) (*NoContent, error) {
	return &NoContent{
		Status: http.StatusNoContent,
	}, nil
}

func (c *V1BetaProbesController) Ready(ctx context.Context, req *V1BetaProbesStartupRequest) (*NoContent, error) {
	return &NoContent{
		Status: http.StatusNoContent,
	}, nil
}
