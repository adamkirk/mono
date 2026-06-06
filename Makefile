
DC=docker compose --env-file ./.dockerenv -f ./docker-compose.yml -p panoptes
TTY_FLAGS=-it

ifeq ($(NO_TTY),1)
TTY_FLAGS=-T
endif

.PHONY: npm \
	setup-env setup-certs-mkcert setup-hooks setup-hosts-txeh setup \
	up down restart dc build \
	tail-% exec-% \
	api-tidy api-test hook-api-fmtcheck api-fmt \
	api-db-migrate api-db-migrate-down api-db-wipe api-db-hard-reset \
	api-sqlc-gen api-export-river-migrations \

setup-env:
	@[ -f .dockerenv ] || cp .dockerenv.example .dockerenv

setup-certs-mkcert:
	@which mkcert > /dev/null || { echo "Error: mkcert not installed. Run: brew install mkcert"; exit 1; }
	@mkcert -install
	@mkdir -p .local/traefik/certs
	@cd .local/traefik/certs && mkcert "*.panoptestest.dev"

setup-hooks:
	git config core.hooksPath .githooks

setup-hosts-txeh:
	txeh add 127.0.0.1 \
		panoptestest.dev api.panoptestest.dev \
		traefik.panoptestest.dev pgadmin.panoptestest.dev \
		--comment "panoptes-dev"

setup: setup-env setup-certs-mkcert setup-hooks

npm:
	@[ $$(node -v | tr -d v | cut -d. -f1) -ge 25 ] || { echo "Error: node $$(node -v) < required v25"; exit 1; }
	@[ $$(npm -v | cut -d. -f1) -ge 11 ] || { echo "Error: npm $$(npm -v) < required v11"; exit 1; }
	npm install

up: setup-env
	$(DC) up -d

down:
	$(DC) down --remove-orphans

restart: down up

build:
	$(DC) build

dc:
	@echo "$(DC)"


exec-%:
	@$(DC) exec $(TTY_FLAGS) $* bash

tail-%:
	@$(DC) logs -f $*

api-tidy:
	@$(DC) exec $(TTY_FLAGS) api go mod tidy

api-fmt:
	$(DC) exec $(TTY_FLAGS) api go fmt ./...

api-test:
	@$(DC) exec $(TTY_FLAGS) api ./scripts/test.sh ./...

api-db-migrate:
	@$(DC) run --rm $(TTY_FLAGS) migrate up

api-db-migrate-down:
	@$(DC) run --rm $(TTY_FLAGS) migrate down 1

api-db-wipe:
	@$(DC) run --rm $(TTY_FLAGS) migrate drop

api-db-hard-reset:
	$(MAKE) down
	docker volume rm -f panoptes_postgres
	$(MAKE) up
	# hack until i sort out an actual wait for postgres
	sleep 5
	$(MAKE) api-db-migrate

api-sqlc-gen:
	$(DC) exec $(TTY_FLAGS) api sqlc generate

api-export-river-migrations:
	$(DC) exec $(TTY_FLAGS) api river migrate-get --all --up > ./api/etc/postgres/migrations/0004_river.up.sql
	$(DC) exec $(TTY_FLAGS) api river migrate-get --all --down > ./api/etc/postgres/migrations/0004_river.down.sql

hook-api-fmtcheck:
	$(DC) exec $(TTY_FLAGS) api ./scripts/fmt_check.sh

