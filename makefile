.PHONY: make life easy

DOCKER_COMPOSE = docker-compose -f docker-compose.yml

refresh:
	docker exec -it snippetbox go build -o /go/bin/web -v ./cmd/web

bs:
	$(DOCKER_COMPOSE) up --build -d

build:
	$(DOCKER_COMPOSE) build

up-bg:
	$(DOCKER_COMPOSE) up -d

up:
	$(DOCKER_COMPOSE) up

start:
	$(DOCKER_COMPOSE) start

logs:
	$(DOCKER_COMPOSE) logs -f

stop:
	$(DOCKER_COMPOSE) stop

enter:
	docker exec -it snippetbox sh

restart:
	@make -s stop
	@make -s start

destroy:
	$(DOCKER_COMPOSE) down -v


down:
	$(DOCKER_COMPOSE) down
