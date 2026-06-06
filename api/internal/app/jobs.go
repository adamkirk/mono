package app

import "github.com/google/uuid"

// RunDeploymentJob queues the actual helm deploy work for a deployment.
type RunDeploymentJob struct {
	DeploymentID uuid.UUID `json:"deployment_id"`
}

func (RunDeploymentJob) Kind() string { return "deployments.run" }
