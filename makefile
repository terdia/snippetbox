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

stop:
	$(DOCKER_COMPOSE) stop

shell:
	docker exec -it snippetbox sh

db-shell:
	docker exec -it database sh

logs:
	$(DOCKER_COMPOSE) logs -f

down:
	$(DOCKER_COMPOSE) down

restart:
	@make -s stop
	@make -s start

destroy:
	$(DOCKER_COMPOSE) down -v

