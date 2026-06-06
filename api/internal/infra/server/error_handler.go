package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/danielgtaylor/huma/v2"
)

func ErrorHandler[Req any, Resp any](handler func(context.Context, *Req) (*Resp, error), method string) func(ctx context.Context, req *Req) (*Resp, error) {
	return func(ctx context.Context, req *Req) (*Resp, error) {
		resp, err := handler(ctx, req)

		if err == nil {
			return resp, nil
		}

		if conflictErr, ok := errors.AsType[common.ErrConflict](err); ok {
			return resp, huma.Error409Conflict(conflictErr.Error())
		}

		switch e := err.(type) {
		case common.ValidationError:
			return resp, buildValidationError(req, method, e)

		case common.ErrUnauthorised:
			return resp, huma.Error401Unauthorized(e.Message)

		case common.ErrNotFound:
			return resp, huma.Error404NotFound(e.Message)

		default:
			// TODO: outside dev, return a generic errors message instead of the
			// actual error as it appears in the response
			return resp, err
		}
	}
}

type MapsErrorKeys interface {
	MapErrorKey(string) string
}

type errorKeyMapper func(string) string

func defaultKeyMapper(target string) string {
	// Do nothing really...
	return target
}

func buildValidationError(req any, method string, err common.ValidationError) *huma.ErrorModel {
	var fieldMapper errorKeyMapper = defaultKeyMapper

	if mapper, ok := req.(MapsErrorKeys); ok {
		fieldMapper = mapper.MapErrorKey
	}

	var statusCode = http.StatusUnprocessableEntity

	errModel := &huma.ErrorModel{
		Title:  http.StatusText(statusCode),
		Detail: "validation failed",
		Status: statusCode,
	}

	for _, fldError := range err.FieldErrors {
		errModel.Add(&huma.ErrorDetail{
			Message:  fldError.Error,
			Location: fieldMapper(fldError.Key),
			Value:    fldError.Value,
		})
	}
	return errModel
}
