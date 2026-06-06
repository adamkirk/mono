package common

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

type Environment struct {
	ID   uuid.UUID
	Name string
}

type EnvironmentComponent struct {
	ID            uuid.UUID
	EnvironmentID uuid.UUID
	Name          string
	ChartName     string
	ChartVersion  string
	ChartRegistry string
}

func IsValidSlug(name string) bool {
	return slugRegex.MatchString(name)
}

type DeploymentStatus string

const DeploymentStatusRequested = "requested"
const DeploymentStatusPending = "pending"
const DeploymentStatusRunning = "running"
const DeploymentStatusCompleted = "completed"
const DeploymentStatusFailed = "failed"

type Deployment struct {
	ID                     uuid.UUID
	CreatedAt              time.Time
	Status                 DeploymentStatus
	EnvironmentID          uuid.UUID
	EnvironmentComponentID uuid.UUID
}
