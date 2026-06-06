CREATE TABLE IF NOT EXISTS "environments"(
   "id" UUID PRIMARY KEY,
   "name" TEXT NOT NULL,
   CONSTRAINT "environments_name_unique" UNIQUE ("name")
);

COMMENT ON COLUMN "environments"."id" IS 'Primary key and unique ID for an environment in the system.';
COMMENT ON COLUMN "environments"."name" IS 'Preferred name to display when referencing this user, free text could be a real name, partial name, nickname etc.';
COMMENT ON CONSTRAINT "environments_name_unique" ON "environments" IS 'We should never have more than one environment with the same name.';
