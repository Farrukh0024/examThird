package main

import (
	"context"
	"fmt"

	"kitab/api"
	_ "kitab/api/docs"
	"kitab/config"
	"kitab/service"
	"kitab/storage/postgres"
)

func main() {
	cfg := config.Load()

	pgStore, err := postgres.New(context.Background(), cfg)
	if err != nil {
		fmt.Println("error while connecting to db", err)
		return
	}
	defer pgStore.Close()

	services := service.New(pgStore)

	server := api.New(services)

	fmt.Println("Service is running on", "port", 8080)
	if err = server.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
