CREATE TABLE IF NOT EXISTS "environment_components"(
   "id" UUID PRIMARY KEY,
   "environment_id" UUID NOT NULL,
   "name" TEXT NOT NULL,
   "chart_name" TEXT NOT NULL,
   "chart_version" TEXT NOT NULL,
   "chart_registry" TEXT NOT NULL,
   CONSTRAINT "environment_components_environment_id_fk" FOREIGN KEY ("environment_id") REFERENCES "environments"("id"),
   CONSTRAINT "environment_components_environment_name_unique" UNIQUE ("environment_id", "name")
);

COMMENT ON COLUMN "environment_components"."id" IS 'Primary key and unique ID for a component in the system.';
COMMENT ON COLUMN "environment_components"."environment_id" IS 'The environment this component belongs to.';
COMMENT ON COLUMN "environment_components"."name" IS 'Name of the component, unique within an environment.';
COMMENT ON COLUMN "environment_components"."chart_name" IS 'Name of the Helm chart.';
COMMENT ON COLUMN "environment_components"."chart_version" IS 'Version of the Helm chart.';
COMMENT ON COLUMN "environment_components"."chart_registry" IS 'Registry where the Helm chart is hosted.';
COMMENT ON CONSTRAINT "environment_components_environment_name_unique" ON "environment_components" IS 'Component names must be unique within an environment.';
