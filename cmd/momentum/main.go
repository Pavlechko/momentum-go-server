package main

import (
	"momentum-go-server/internal/handlers"
	"momentum-go-server/internal/routes"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"
)

func main() {
	db, err := store.ConnectDB()
	if err != nil {
		utils.ErrorLogger.Println("Error connectiont to database:", err)
	}
	defer db.Close()

	service := services.CreateService(db)
	handler := handlers.CreateHandler(service)

	server := routes.NewAPIServer(":8080", handler)
	server.Run()
}
