package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/google/uuid"
)

type DeploymentsHandler struct {
	l                               *slog.Logger
	environmentsRepository          environmentsRepository
	environmentComponentsRepository environmentComponentsRepository
	deploymentsRepository           deploymentsRepository
}

type CreateDeploymentDTO struct {
	EnvironmentName string
	ComponentName   string
}

func (dto CreateDeploymentDTO) Validate() error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.EnvironmentName) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "EnvironmentName",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.EnvironmentName,
		})
	}

	if !common.IsValidSlug(dto.ComponentName) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "ComponentName",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.ComponentName,
		})
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *DeploymentsHandler) Create(dto CreateDeploymentDTO) (*common.Deployment, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{Message: "the environment was not found"}
	}

	component, err := h.environmentComponentsRepository.ByEnvironmentAndName(env.ID, dto.ComponentName)
	if err != nil {
		return nil, err
	}

	if component == nil {
		return nil, common.ErrNotFound{Message: "the component was not found"}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	d := &common.Deployment{
		ID:                     id,
		CreatedAt:              time.Now(),
		Status:                 common.DeploymentStatusRequested,
		EnvironmentID:          env.ID,
		EnvironmentComponentID: component.ID,
	}

	return d, h.deploymentsRepository.Save(d,
		common.WithQueuedJob(func(q common.JobEnQueuer) error {
			return q.Enqueue(context.Background(), RunDeploymentJob{DeploymentID: d.ID})
		}),
	)
}

type GetDeploymentDTO struct {
	EnvironmentName string
	ComponentName   string
	DeploymentID    uuid.UUID
}

func (h *DeploymentsHandler) Get(dto GetDeploymentDTO) (*common.Deployment, error) {
	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{Message: "the environment was not found"}
	}

	component, err := h.environmentComponentsRepository.ByEnvironmentAndName(env.ID, dto.ComponentName)
	if err != nil {
		return nil, err
	}

	if component == nil {
		return nil, common.ErrNotFound{Message: "the component was not found"}
	}

	deployment, err := h.deploymentsRepository.ByID(dto.DeploymentID)
	if err != nil {
		return nil, err
	}

	if deployment == nil || deployment.EnvironmentComponentID != component.ID {
		return nil, common.ErrNotFound{}
	}

	return deployment, nil
}

func NewDeploymentsHandler(
	l *slog.Logger,
	environmentsRepository environmentsRepository,
	environmentComponentsRepository environmentComponentsRepository,
	deploymentsRepository deploymentsRepository,
) *DeploymentsHandler {
	return &DeploymentsHandler{
		l:                               l,
		environmentsRepository:          environmentsRepository,
		environmentComponentsRepository: environmentComponentsRepository,
		deploymentsRepository:           deploymentsRepository,
	}
}
