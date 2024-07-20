# my keyboard is 60% i hate this kind of keyboard
# why i bought this shit of keyboard? OMG

DOCKER_COMPOSE = docker-compose
BUILD = $(DOCKER_COMPOSE) build
UP = $(DOCKER_COMPOSE) up -d
DOWN = $(DOCKER_COMPOSE) down
LOGS = $(DOCKER_COMPOSE) logs -f

#dont need this but ok
.PHONY: build up down logs

build:
	@$(BUILD)

up: build
	@$(UP)

down:
	@$(DOWN)

logs:
	@$(LOGS)

start: up logs
