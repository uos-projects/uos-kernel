INFRA_COMPOSE=infra/docker-compose.yml

.PHONY: infra-up infra-down

infra-up:
	docker compose -f $(INFRA_COMPOSE) up -d

infra-down:
	docker compose -f $(INFRA_COMPOSE) down -v


