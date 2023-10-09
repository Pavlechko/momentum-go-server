package main

import (
	"momentum-go-server/internal/routes"
	"momentum-go-server/internal/store"
)

func main() {
	store.ContentDB()
	// store.CreateUser()
	// store.GetUser()
	routes.SetupRoutes()
	// time.Sleep(10 * time.Second)
}
