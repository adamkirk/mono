package app

import (
	"log/slog"

	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/google/uuid"
)

type EnvironmentComponentsHandler struct {
	l                               *slog.Logger
	environmentsRepository          environmentsRepository
	environmentComponentsRepository environmentComponentsRepository
}
type CreateEnvironmentComponentDTO struct {
	EnvironmentName string
	Name            string
	ChartName       string
	ChartVersion    string
	ChartRegistry   string
}

func (dto CreateEnvironmentComponentDTO) Validate(repo environmentComponentsRepository, environmentID uuid.UUID) error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})
	} else {
		existing, err := repo.ByEnvironmentAndName(environmentID, dto.Name)
		if err != nil {
			return err
		}

		if existing != nil {
			fldErrors = append(fldErrors, common.FieldError{
				Key:   "Name",
				Error: "a component with this name already exists in this environment",
				Value: dto.Name,
			})
		}
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *EnvironmentComponentsHandler) Create(dto CreateEnvironmentComponentDTO) (*common.EnvironmentComponent, error) {
	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{}
	}

	if err := dto.Validate(h.environmentComponentsRepository, env.ID); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	c := &common.EnvironmentComponent{
		ID:            id,
		EnvironmentID: env.ID,
		Name:          dto.Name,
		ChartName:     dto.ChartName,
		ChartVersion:  dto.ChartVersion,
		ChartRegistry: dto.ChartRegistry,
	}

	return c, h.environmentComponentsRepository.Save(c)
}

type GetEnvironmentComponentDTO struct {
	EnvironmentName string
	Name            string
}

func (dto GetEnvironmentComponentDTO) Validate() error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.EnvironmentName) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "EnvironmentName",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.EnvironmentName,
		})
	}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *EnvironmentComponentsHandler) Get(dto GetEnvironmentComponentDTO) (*common.EnvironmentComponent, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{
			Message: "the environment was not found",
		}
	}

	return h.environmentComponentsRepository.ByEnvironmentAndName(env.ID, dto.Name)
}

type UpdateEnvironmentComponentDTO struct {
	EnvironmentName string
	CurrentName     string
	Name            *string
	ChartName       *string
	ChartVersion    *string
	ChartRegistry   *string
}

func (dto UpdateEnvironmentComponentDTO) Validate(repo environmentComponentsRepository, environmentID uuid.UUID) error {
	if dto.Name == nil && dto.ChartName == nil && dto.ChartVersion == nil && dto.ChartRegistry == nil {
		return common.ValidationError{
			FieldErrors: []common.FieldError{
				{
					Key:   "Body",
					Error: "at least one field must be provided",
				},
			},
		}
	}

	fldErrors := []common.FieldError{}

	if dto.Name != nil {
		if !common.IsValidSlug(*dto.Name) {
			fldErrors = append(fldErrors, common.FieldError{
				Key:   "Name",
				Error: "must contain alphanumeric or hyphen characters only",
				Value: *dto.Name,
			})
		} else {
			existing, err := repo.ByEnvironmentAndName(environmentID, *dto.Name)
			if err != nil {
				return err
			}

			if existing != nil {
				msg := "a component with this name already exists in this environment"
				if *dto.Name == dto.CurrentName {
					msg = "the names are the same"
				}
				fldErrors = append(fldErrors, common.FieldError{
					Key:   "Name",
					Error: msg,
					Value: *dto.Name,
				})
			}
		}
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *EnvironmentComponentsHandler) Update(dto UpdateEnvironmentComponentDTO) (*common.EnvironmentComponent, error) {
	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{Message: "the environment was not found"}
	}

	c, err := h.environmentComponentsRepository.ByEnvironmentAndName(env.ID, dto.CurrentName)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, common.ErrNotFound{}
	}

	if err := dto.Validate(h.environmentComponentsRepository, env.ID); err != nil {
		return nil, err
	}

	if dto.Name != nil {
		c.Name = *dto.Name
	}
	if dto.ChartName != nil {
		c.ChartName = *dto.ChartName
	}
	if dto.ChartVersion != nil {
		c.ChartVersion = *dto.ChartVersion
	}
	if dto.ChartRegistry != nil {
		c.ChartRegistry = *dto.ChartRegistry
	}

	return c, h.environmentComponentsRepository.Save(c)
}

type DeleteEnvironmentComponentDTO struct {
	EnvironmentName string
	Name            string
}

func (dto DeleteEnvironmentComponentDTO) Validate() error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.EnvironmentName) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "EnvironmentName",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.EnvironmentName,
		})
	}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *EnvironmentComponentsHandler) Delete(dto DeleteEnvironmentComponentDTO) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return err
	}

	if env == nil {
		return common.ErrNotFound{Message: "the environment was not found"}
	}

	c, err := h.environmentComponentsRepository.ByEnvironmentAndName(env.ID, dto.Name)
	if err != nil {
		return err
	}

	if c == nil {
		return common.ErrNotFound{}
	}

	return h.environmentComponentsRepository.Delete(c)
}

type ListEnvironmentComponentsDTO struct {
	EnvironmentName string
	Page            int
	PerPage         int
}

func (dto ListEnvironmentComponentsDTO) Validate(totalPages int) error {
	fldErrors := []common.FieldError{}

	if dto.PerPage < 1 || dto.PerPage > 100 {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "PerPage",
			Error: "must be between 1 and 100",
			Value: dto.PerPage,
		})
	}

	// Page 1 is always valid even when there are no records.
	if dto.Page > 1 && dto.Page > totalPages {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Page",
			Error: "page is out of bounds",
			Value: dto.Page,
		})
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

type ListEnvironmentComponentsResult struct {
	Components []*common.EnvironmentComponent
	Total      int
	TotalPages int
	Page       int
	PerPage    int
}

func (h *EnvironmentComponentsHandler) List(dto ListEnvironmentComponentsDTO) (*ListEnvironmentComponentsResult, error) {
	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{}
	}

	total, err := h.environmentComponentsRepository.CountByEnvironment(env.ID)
	if err != nil {
		return nil, err
	}

	totalPages := total / dto.PerPage
	if total%dto.PerPage != 0 {
		totalPages++
	}

	if err := dto.Validate(totalPages); err != nil {
		return nil, err
	}

	offset := (dto.Page - 1) * dto.PerPage
	components, err := h.environmentComponentsRepository.ListByEnvironment(env.ID, dto.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListEnvironmentComponentsResult{
		Components: components,
		Total:      total,
		TotalPages: totalPages,
		Page:       dto.Page,
		PerPage:    dto.PerPage,
	}, nil
}

func NewEnvironmentComponentsHandler(
	l *slog.Logger,
	environmentsRepository environmentsRepository,
	environmentComponentsRepository environmentComponentsRepository,
) *EnvironmentComponentsHandler {
	return &EnvironmentComponentsHandler{
		l:                               l,
		environmentsRepository:          environmentsRepository,
		environmentComponentsRepository: environmentComponentsRepository,
	}
}
