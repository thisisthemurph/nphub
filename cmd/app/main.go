package main

import (
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	appbuilder "nphud/cmd/app/setup/application_builder"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	builder := appbuilder.NewApplicationBuilder()
	builder.AddCore()
	builder.AddInfrastructure()
	builder.AddServices()
	builder.AddRoutes()

	app := builder.Build()

	if err := app.MigrateDatabase(); err != nil {
		log.Fatal(err)
	}

	if err := app.ConfigMediator(); err != nil {
		log.Fatal(err)
	}

	app.MapEndpoints()
	app.Run()
}
