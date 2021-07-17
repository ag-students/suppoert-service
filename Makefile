.PHONY: all

include .env
export

PROJ_PATH := ${CURDIR}
DOCKER_PATH := ${PROJ_PATH}/docker

APP=communication
MIGRATION_TOOL=goose
MIGRATIONS_DIR=./db/migrations

BASIC_IMAGE=dep
IMAGE_POSTFIX=-image

build:
	go build -o .bin/communication cmd/communication/main.go
	go build -o .bin/docsgenerator cmd/docsgenerator/main.go
	chmod ugo+x .bin/communication
	chmod ugo+x .bin/docsgenerator

build-docker:
	sudo rm -rf .database/
	docker build -t ${BASIC_IMAGE} -f ${DOCKER_PATH}/builder.Dockerfile.dev .
	docker build -t communication${IMAGE_POSTFIX} -f ${DOCKER_PATH}/communication.Dockerfile.dev .
	docker build -t docsgenerator${IMAGE_POSTFIX} -f ${DOCKER_PATH}/docsgenerator.Dockerfile.dev .
app-setup-and-up: build-docker app-up

app-up: build
	docker-compose up

all: app-setup-and-up

app-bash:
	docker-compose run --rm --no-deps --name communication-service ${APP} bash

app-up-local: build
	./.bin/communication

db-bash:
	docker-compose run --rm --no-deps --name communication-db db ash

goose-init:
	go build -o .bin/goose cmd/${MIGRATION_TOOL}/main.go
	chmod ugo+x .bin/${MIGRATION_TOOL}

db-up:
	docker-compose run --rm --no-deps --name communication-db db ash

db-migration-create: goose-init
	if [ -z ${lang} ] ; \
	then \
		goose -dir=${MIGRATIONS_DIR} create ${name} sql ; \
	else \
	  	goose -dir=${MIGRATIONS_DIR} create ${name} ${lang} ; \
	fi ;

db-migrate-status: goose-init
	docker-compose run --rm communication .bin/goose -dir ${MIGRATIONS_DIR} postgres \
		"user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL}" status

db-migrate-up: goose-init
	docker-compose run --rm communication .bin/goose -dir ${MIGRATIONS_DIR} postgres \
        "user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL}" up

db-migrate-down: goose-init
	docker-compose run --rm communication .bin/goose -dir ${MIGRATIONS_DIR} postgres \
        "user=${POSTGRES_USER} host=${POSTGRES_HOST} port=${POSTGRES_PORT} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL}" down

test:
	gotest -v ./...

packages-tidy:
	go mod tidy

# Kakafka

kafka-bash:
	docker exec -it kafka bash

create-topic:
	docker exec -it kafka /opt/kafka_2.13-2.7.0/bin/kafka-topics.sh \
		--create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic ${name}