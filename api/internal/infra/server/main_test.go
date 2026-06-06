package server_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/adamkirk/panoptes/api/internal/infra/server"
	"github.com/danielgtaylor/huma/v2"
)

var BlahVersion server.ApiVersion = "blah"

type DummyController struct{}

func (c *DummyController) RegisterRoutes(v server.ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.dummy.healthz", string(v)),
		Method:        http.MethodGet,
		Path:          "/healthz",
		Summary:       "Check if the app is started up",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Healthz",
		},
		Metadata: map[string]any{},
	}, c.Healthz)
}

type HealthzRequest struct{}

func (c *DummyController) Healthz(ctx context.Context, req *HealthzRequest) (*server.NoContent, error) {
	return &server.NoContent{
		Status: http.StatusNoContent,
	}, nil
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:noctx
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}

	return fmt.Errorf("server did not become ready within %s", timeout)
}

func TestAPI(t *testing.T) {
	port := freePort(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := server.New(
		port,
		logger,
		server.WithApiVersionGroup(server.ApiVersionGroup{
			Version: server.ApiVersionV1Beta,
			Controllers: []server.Controller{
				&DummyController{},
			},
		}),
		server.WithApiVersionGroup(server.ApiVersionGroup{
			Version: BlahVersion,
			Controllers: []server.Controller{
				&DummyController{},
			},
		}),
	)

	go s.Start() //nolint:errcheck

	base := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForServer(base+"/api/v1beta/healthz", 2*time.Second); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(base + "/api/v1beta/healthz") //nolint:noctx

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}

	resp, err = http.Get(base + "/api/blah/healthz") //nolint:noctx

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}
