.PHONY: all

include .env
export

PROJ_PATH := ${CURDIR}
DOCKER_PATH := ${PROJ_PATH}/docker

APP=communication
MIGRATION_TOOL=goose
MIGRATIONS_DIR=migrations

BASIC_IMAGE=default
IMAGE_POSTFIX=-image

build-app:
	go build -o .bin/${APP} cmd/${APP}/main.go
	chmod ugo+x .bin/${APP}

build-docker:
	docker build -t ${BASIC_IMAGE} -f ${DOCKER_PATH}/builder.Dockerfile .
	docker build -t ${APP}${IMAGE_POSTFIX} -f ${DOCKER_PATH}/${APP}.Dockerfile .

app-setup-and-up: build-docker app-up

app-up: build-app
	docker-compose up

all: app-setup-and-up

app-bash:
	docker-compose run --rm --no-deps --name communication-service ${APP} bash

app-up-local: build-app
	./.bin/communication

db-bash:
	docker-compose run --rm --no-deps --name communication-db db ash

goose-init:
	go build -o .bin/goose cmd/${MIGRATION_TOOL}/main.go
	chmod ugo+x .bin/${MIGRATION_TOOL}

db-up:
	docker-compose run --rm --no-deps --name communication-db db ash

db-migrate-status: goose-init
	docker-compose run --rm events-server .bin/goose -dir ${MIGRATIONS_DIR} postgres \
		"user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_DBPASSWORD} dbname=${POSTGRES_DBNAME} sslmode=${POSTGRES_SSL}" status

db-migrate-up: goose-init
	docker-compose run --rm events-server .bin/goose -dir ${MIGRATIONS_DIR} postgres \
        "user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_DBPASSWORD} dbname=${POSTGRES_DBNAME} sslmode=${POSTGRES_SSL}" up

db-migrate-down: goose-init
	docker-compose run --rm events-server .bin/goose -dir ${MIGRATIONS_DIR} postgres \
        "user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_DBPASSWORD} dbname=${POSTGRES_DBNAME} sslmode=${POSTGRES_SSL}" down