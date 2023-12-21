package main

import (
	"context"
	"momentum-go-server/internal/handlers"
	"momentum-go-server/internal/routes"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/store/redis"
	"momentum-go-server/internal/utils"
)

func main() {
	var ctx = context.Background()

	db, err := store.ConnectDB()
	if err != nil {
		utils.ErrorLogger.Println("Error connectiont to database:", err)
	}
	defer db.Close()

	redisClient := redis.NewRedisClient("localhost:6379", 0, ctx)

	service := services.CreateService(db, *redisClient)
	handler := handlers.CreateHandler(service)

	server := routes.NewAPIServer(":8080", handler)
	server.Run()
}
