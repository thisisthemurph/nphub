# Build, run, and test

test:
	@go test -v ./...

app:
	@templ generate
	@go run ./cmd/app

bg:
	@go run ./cmd/background

# Tailwind

get-tw:
	@rm -f tailwindcss
	@curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
	@chmod +x tailwindcss-linux-x64
	@mv tailwindcss-linux-x64 tailwindcss

tw:
	@./tailwindcss -i ./tailwind.css -o ./internal/app/shared/static/css/app.css --watch

tailwind-once:
	@./tailwindcss -i ./tailwind.css -o ./internal/app/shared/static/css/app.css

tw-refresh:
	@truncate -s 0 internal/view/static/app.css
	@make tw

# Database migrations

up:
	@go run cmd/migrate/main.go up

down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))