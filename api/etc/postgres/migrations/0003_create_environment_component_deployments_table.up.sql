CREATE TABLE IF NOT EXISTS "environment_component_deployments"(
   "id" UUID PRIMARY KEY,
   "environment_id" UUID NOT NULL,
   "environment_component_id" UUID NOT NULL,
   "created_at" TIMESTAMPTZ NOT NULL,
   "status" TEXT NOT NULL,
   CONSTRAINT "environment_component_deployments_environment_id_fk" FOREIGN KEY ("environment_id") REFERENCES "environments"("id"),
   CONSTRAINT "environment_component_deployments_environment_component_id_fk" FOREIGN KEY ("environment_component_id") REFERENCES "environment_components"("id")
);

COMMENT ON COLUMN "environment_component_deployments"."id" IS 'Primary key and unique ID for a deployment in the system.';
COMMENT ON COLUMN "environment_component_deployments"."environment_id" IS 'The environment that this deploy was part of.';
COMMENT ON COLUMN "environment_component_deployments"."environment_component_id" IS 'The environment component that was deployed.';
COMMENT ON COLUMN "environment_component_deployments"."created_at" IS 'When the deployment was created.';
COMMENT ON COLUMN "environment_component_deployments"."status" IS 'The status of the deployment.';
