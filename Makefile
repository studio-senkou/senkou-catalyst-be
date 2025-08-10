# Export only Discord variables if .env exists
ifneq (,$(wildcard ./.env))
DISCORD_BOT_URL := $(shell grep '^DISCORD_BOT_URL=' .env | cut -d '=' -f2-)
DISCORD_CHANNEL_ID := $(shell grep '^DISCORD_CHANNEL_ID=' .env | cut -d '=' -f2-)
export DISCORD_BOT_URL
export DISCORD_CHANNEL_ID
endif

.PHONY=auth-secret rebuild rebuild-dev wire dev-up dev-down dev-logs prod-up prod-down prod-logs seed clean dev-status prod-status list-all swagger test-discord

auth-secret:
	@echo "" >> .env
	@echo "AUTH_SECRET=\"$(shell openssl rand -base64 32 | tr -d '\n')\"" >> .env
	@echo "AUTH_SECRET generated and saved to .env"

rebuild:
	@echo "🚀 Starting production rebuild..."
	@$(call send_discord_notification,🚀 Production rebuild started for Senkou Catalyst BE)
	@mv .env .env.temp 2>/dev/null || true
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod down --remove-orphans
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod build --no-cache
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod up -d
	@mv .env.temp .env 2>/dev/null || true
	@docker image prune -f
	@echo "✅ Production rebuild completed successfully!"
	@$(call send_discord_notification,✅ Production rebuild completed successfully for Senkou Catalyst BE)

rebuild-stage:
	@echo "🚀 Starting staging rebuild..."
	@$(call send_discord_notification,🚀 Staging rebuild started for Senkou Catalyst BE)
	@mv .env .env.temp 2>/dev/null || true
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging down --remove-orphans
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging build --no-cache
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging up -d
	@mv .env.temp .env 2>/dev/null || true
	@docker image prune -f
	@echo "✅ Staging rebuild completed successfully!"
	@$(call send_discord_notification,✅ Staging rebuild completed successfully for Senkou Catalyst BE)

wire:
	@cd container && wire
	@echo "Wire dependency injection generated successfully."

stage-up:
	@mv .env .env.temp 2>/dev/null || true
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging up -d
	@mv .env.temp .env 2>/dev/null || true
	@echo "Staging environment started"

stage-down:
	@mv .env .env.temp 2>/dev/null || true
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging down
	@mv .env.temp .env 2>/dev/null || true
	@echo "Staging environment stopped"

stage-logs:
	@docker compose -f docker-compose-staging.yml --env-file .env.staging -p senkou-catalyst-staging logs -f

prod-up:
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod up -d
	@echo "Production environment started"

prod-down:
	@docker compose -f docker-compose.yml --env.file .env.production -p senkou-catalyst-prod down
	@echo "Production environment stopped"

prod-logs:
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod logs -f

seed:
	@go run cmd/seed.go
	@echo "Database seeded successfully"

clean:
	@docker system prune -f
	@docker image prune -f
	@echo "Docker cleanup completed"

dev-status:
	@echo "=== Development Environment Status ==="
	@docker compose -f docker-compose-dev.yml --env-file .env.development -p senkou-catalyst-dev ps
	@echo ""
	@echo "=== Development Volumes ==="
	@docker volume ls | grep dev || echo "No dev volumes found"

prod-status:
	@echo "=== Production Environment Status ==="
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-catalyst-prod ps
	@echo ""
	@echo "=== Production Volumes ==="
	@docker volume ls | grep prod || echo "No prod volumes found"

list-all:
	@echo "=== All Containers ==="
	@docker ps -a --filter name=senkou-catalyst --format "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}"
	@echo ""
	@echo "=== All Networks ==="
	@docker network ls | grep catalyst || echo "No catalyst networks found"
	@echo ""
	@echo "=== All Volumes ==="
	@docker volume ls | grep catalyst || echo "No catalyst volumes found"

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -o ./docs --parseDependency --parseInternal
	@echo "Swagger documentation generated successfully!"

test-discord:
	@echo "Testing Discord notification..."
	@$(call send_discord_notification,🧪 Test notification from Senkou Catalyst BE Makefile)
	@echo "Discord notification test sent!"

define send_discord_notification
	@if [ -n "$(DISCORD_BOT_URL)" ] && [ -n "$(DISCORD_CHANNEL_ID)" ]; then \
		curl --silent --request POST \
		--location "$(DISCORD_BOT_URL)/send-log" \
		--header "Content-Type: application/json" \
		--data "{\"message\": \"$(1)\", \"channelId\": \"$(DISCORD_CHANNEL_ID)\"}" \
		> /dev/null 2>&1 || true; \
	fi
endef