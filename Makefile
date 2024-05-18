test:
	@go test -v ./...

app:
	@templ generate
	@go run ./cmd/app

bg:
	@go run ./cmd/background

up:
	@go run cmd/migrate/main.go up

down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))