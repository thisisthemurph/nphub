package main

import (
	"fmt"
	"nphud/internal/repository"
	"nphud/pkg/store"
)

func main() {
	database, err := store.GetOrCreate()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	repo := repository.NewGameRepository(database)
	games, err := repo.List()
	if err != nil {
		panic(err)
	}

	for _, g := range games {
		fmt.Printf("Game ID: %d\n", g.ID)
	}
}
