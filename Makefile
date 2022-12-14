POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=medium_user_service_db

CURRENT_DIR=$(shell pwd)

-include .env
  
DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"




run:
	go run cmd/main.go

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate_file:
	migrate create -ext sql -dir migrations/ -seq alter_some_table

local-up:
	docker compose --env-file ./.env.docker up -d

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}

pull-sub-module:
	git submodule update --init --recursive

update-sub-module:
	git submodule update --remote --merge


.PHONY: run migrateup migratedown local-up proto-gen pull-sub-module update-sub-module