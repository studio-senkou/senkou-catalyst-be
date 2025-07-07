.PHONY=auth-secret rebuild

auth-secret:
	@echo "" >> .env
	@echo "AUTH_SECRET=\"$(shell openssl rand -base64 32 | tr -d '\n')\"" >> .env
	@echo "AUTH_SECRET generated and saved to .env"

rebuild:
	@docker compose down
	@docker compose build --no-cache
	@docker compose up -d
	@docker image prune -f