package app

import (
	"log/slog"

	"github.com/adamkirk/panoptes/api/internal/common"
	"github.com/google/uuid"
)

type EnvironmentsHandler struct {
	l                      *slog.Logger
	environmentsRepository environmentsRepository
}

type CreateEnvironmentDTO struct {
	Name string
}

func (dto CreateEnvironmentDTO) Validate(repo environmentsRepository) error {
	found, err := repo.ByName(dto.Name)

	if err != nil {
		return err
	}

	fldErrors := []common.FieldError{}

	if found != nil {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "an environment with this name already exists.",
			Value: dto.Name,
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
		return common.ValidationError{
			FieldErrors: fldErrors,
		}
	}

	return nil
}

func (h *EnvironmentsHandler) Create(dto CreateEnvironmentDTO) (*common.Environment, error) {
	if err := dto.Validate(h.environmentsRepository); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	env := &common.Environment{
		ID:   id,
		Name: dto.Name,
	}

	return env, h.environmentsRepository.Save(env)
}

type GetEnvironmentDTO struct {
	Name string
}

func (dto GetEnvironmentDTO) Validate() error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})

	}
	if len(fldErrors) > 0 {
		return common.ValidationError{
			FieldErrors: fldErrors,
		}
	}

	return nil
}

func (h *EnvironmentsHandler) Get(dto GetEnvironmentDTO) (*common.Environment, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return h.environmentsRepository.ByName(dto.Name)
}

type UpdateEnvironmentDTO struct {
	CurrentName string
	Name        string
}

func (dto UpdateEnvironmentDTO) Validate(repo environmentsRepository) error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})
	} else {
		existing, err := repo.ByName(dto.Name)
		if err != nil {
			return err
		}

		if existing != nil {
			msg := "an environment with this name already exists"

			if existing.Name == dto.Name {
				msg = "the names are the same"
			}

			fldErrors = append(fldErrors, common.FieldError{
				Key:   "Name",
				Error: msg,
				Value: dto.Name,
			})
		}
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

func (h *EnvironmentsHandler) Update(dto UpdateEnvironmentDTO) (*common.Environment, error) {
	env, err := h.environmentsRepository.ByName(dto.CurrentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{}
	}

	if err := dto.Validate(h.environmentsRepository); err != nil {
		return nil, err
	}

	env.Name = dto.Name

	if err := h.environmentsRepository.Save(env); err != nil {
		return nil, err
	}

	return env, nil
}

type DeleteEnvironmentDTO struct {
	Name string
}

func (dto DeleteEnvironmentDTO) Validate() error {
	fldErrors := []common.FieldError{}

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

func (h *EnvironmentsHandler) Delete(dto DeleteEnvironmentDTO) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	env, err := h.environmentsRepository.ByName(dto.Name)
	if err != nil {
		return err
	}

	if env == nil {
		return common.ErrNotFound{}
	}

	return h.environmentsRepository.Delete(env)
}

type ListEnvironmentsDTO struct {
	Page    int
	PerPage int
}

func (dto ListEnvironmentsDTO) Validate(totalPages int) error {
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

type ListEnvironmentsResult struct {
	Environments []*common.Environment
	Total        int
	TotalPages   int
	Page         int
	PerPage      int
}

func (h *EnvironmentsHandler) List(dto ListEnvironmentsDTO) (*ListEnvironmentsResult, error) {
	total, err := h.environmentsRepository.Count()
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
	envs, err := h.environmentsRepository.List(dto.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return &ListEnvironmentsResult{
		Environments: envs,
		Total:        total,
		TotalPages:   totalPages,
		Page:         dto.Page,
		PerPage:      dto.PerPage,
	}, nil
}

func NewEnvironmentsHandler(l *slog.Logger, environmentsRepository environmentsRepository) *EnvironmentsHandler {
	return &EnvironmentsHandler{
		l:                      l,
		environmentsRepository: environmentsRepository,
	}
}
