-- name: GetEnvironmentByName :one
SELECT
    *
FROM environments
WHERE
    name = $1;

-- name: ListEnvironments :many
SELECT *
FROM environments
ORDER BY name ASC
LIMIT $1
OFFSET $2;

-- name: CountEnvironments :one
SELECT COUNT(*) FROM environments;

-- name: DeleteEnvironmentByID :exec
DELETE FROM environments WHERE id = $1;

-- name: UpsertEnvironment :one
INSERT INTO environments (id, name)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: GetEnvironmentComponentByEnvironmentAndName :one
SELECT * FROM environment_components
WHERE environment_id = $1 AND name = $2;

-- name: CountEnvironmentComponents :one
SELECT COUNT(*) FROM environment_components
WHERE environment_id = $1;

-- name: ListEnvironmentComponents :many
SELECT * FROM environment_components
WHERE environment_id = $1
ORDER BY name ASC
LIMIT $2
OFFSET $3;

-- name: DeleteEnvironmentComponentByID :exec
DELETE FROM environment_components WHERE id = $1;

-- name: UpsertEnvironmentComponent :one
INSERT INTO environment_components (id, environment_id, name, chart_name, chart_version, chart_registry)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    chart_name = EXCLUDED.chart_name,
    chart_version = EXCLUDED.chart_version,
    chart_registry = EXCLUDED.chart_registry
RETURNING *;

-- name: GetDeploymentByID :one
SELECT * FROM environment_component_deployments WHERE id = $1;

-- name: UpsertDeployment :one
INSERT INTO environment_component_deployments (id, environment_id, environment_component_id, created_at, status)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status
RETURNING *;