package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const OptDisableNotFound = "DisableNotFound"
const OptDisableAllDefaultResponses = "DisableAllDefaults"
const OptDisableDefaultAuthentication = "DisableAuthentication"

type ApiVersion string

const ApiVersionV1Beta ApiVersion = "v1beta"

type Server struct {
	port             int
	echo             *echo.Echo
	api              huma.API
	accessLogger     *slog.Logger
	apiVersionGroups []ApiVersionGroup
}

type Controller interface {
	RegisterRoutes(v ApiVersion, g huma.API)
}

type ApiVersionGroup struct {
	Version     ApiVersion
	Controllers []Controller
}

type ServerOpt func(s *Server)

// versionToSchemaPrefix converts an ApiVersion like "v1beta" to a Go type
// prefix like "V1Beta" so it can be stripped from generated schema names.
// Spat out by claude, not a huge fan, but it seems to work for now...
func versionToSchemaPrefix(v ApiVersion) string {
	s := string(v)
	var b strings.Builder
	prevWasDigit := false
	for i, r := range s {
		switch {
		case i == 0:
			b.WriteRune(unicode.ToUpper(r))
		case prevWasDigit && unicode.IsLetter(r):
			b.WriteRune(unicode.ToUpper(r))
		default:
			b.WriteRune(r)
		}
		prevWasDigit = unicode.IsDigit(r)
	}
	return b.String()
}

func schemaNamerFor(v ApiVersion) func(reflect.Type, string) string {
	prefix := versionToSchemaPrefix(v)
	return func(t reflect.Type, hint string) string {
		return strings.TrimPrefix(huma.DefaultSchemaNamer(t, hint), prefix)
	}
}

func setupHumaMiddlewares(api huma.API) {
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		ctx.SetHeader("X-Operation-Id", ctx.Operation().OperationID)
		next(ctx)
	})
}

func setupEchoMiddlewares(e *echo.Echo, logger *slog.Logger, accessLogger *slog.Logger) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.ErrorContext(c.Request().Context(), "panic recovered",
				slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				slog.String("error", err.Error()),
				slog.String("stack", string(stack)),
			)
			return nil
		},
	}))

	e.Use(middleware.RequestID())

	if accessLogger != nil {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogMethod:    true,
			LogURI:       true,
			LogStatus:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogRequestID: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				level := slog.LevelInfo
				switch {
				case v.Status >= 500:
					level = slog.LevelError
				case v.Status >= 400:
					level = slog.LevelWarn
				}

				accessLogger.LogAttrs(c.Request().Context(), level, "access log",
					slog.String("operation_id", c.Response().Header().Get("X-Operation-Id")),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Duration("latency", v.Latency),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("request_id", v.RequestID),
				)
				return nil
			},
		}))
	}
}

func addValidationErrorResponse(op *huma.Operation) {
	validationStatus := strconv.Itoa(http.StatusUnprocessableEntity)

	if _, ok := op.Responses[validationStatus]; ok {
		return
	}

	op.Responses[validationStatus] = &huma.Response{
		Description: "validation error",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}
func addInternalErrorResponse(op *huma.Operation) {
	internalError := strconv.Itoa(http.StatusInternalServerError)

	if _, ok := op.Responses[internalError]; ok {
		return
	}

	op.Responses[internalError] = &huma.Response{
		Description: "Internal server error",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}

func addNotFoundResponse(op *huma.Operation) {
	var notFoundEnabled = true

	if v, ok := op.Metadata[OptDisableNotFound]; ok {
		if optAsBool, ok := v.(bool); ok {
			notFoundEnabled = !optAsBool
		}
	}

	if !notFoundEnabled {
		return
	}

	notFound := strconv.Itoa(http.StatusNotFound)

	if _, ok := op.Responses[notFound]; ok {
		return
	}

	op.Responses[notFound] = &huma.Response{
		Description: "Resource Not Found",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}

func configureDefaultResponses(api *huma.OpenAPI, op *huma.Operation) {

	if _, ok := op.Responses["default"]; ok {
		// Remove the default as it's an error, but has no status code
		// Maybe there's another way to turn it off
		op.Responses["default"] = nil
	}

	if v, ok := op.Metadata[OptDisableAllDefaultResponses]; ok {
		if optAsBool, ok := v.(bool); ok && optAsBool {
			return
		}
	}

	addValidationErrorResponse(op)
	addInternalErrorResponse(op)
	addNotFoundResponse(op)
}

func setupHumaHooks(api huma.API) {
	api.OpenAPI().OnAddOperation = append(
		api.OpenAPI().OnAddOperation,
		// Note, this should come before default responses, as we may want to use
		// the security requirements to configure extra responses based on whether
		// authentication is required.
		configureDefaultResponses,
	)
}

func WithAccessLogger(l *slog.Logger) ServerOpt {
	return func(s *Server) {
		s.accessLogger = l
	}
}

func WithApiVersionGroup(g ApiVersionGroup) ServerOpt {
	return func(s *Server) {
		s.apiVersionGroups = append(s.apiVersionGroups, g)
	}
}

func New(port int, logger *slog.Logger, opts ...ServerOpt) *Server {
	s := &Server{
		port:             port,
		apiVersionGroups: []ApiVersionGroup{},
	}

	for _, o := range opts {
		o(s)
	}

	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.HidePort = true

	setupEchoMiddlewares(s.echo, logger, s.accessLogger)

	for _, vg := range s.apiVersionGroups {
		prefix := fmt.Sprintf("/api/%s", string(vg.Version))
		g := s.echo.Group(prefix)
		apiCfg := huma.DefaultConfig("Panoptes API", string(vg.Version))
		apiCfg.OpenAPI.Components.Schemas = huma.NewMapRegistry("#/components/schemas/", schemaNamerFor(vg.Version))
		apiCfg.OpenAPI.Servers = []*huma.Server{{URL: prefix}}
		apiCfg.RejectUnknownQueryParameters = true
		hg := humaecho.NewWithGroup(s.echo, g, apiCfg)

		setupHumaMiddlewares(hg)
		setupHumaHooks(hg)

		for _, c := range vg.Controllers {
			c.RegisterRoutes(vg.Version, hg)
		}
	}

	return s
}

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
