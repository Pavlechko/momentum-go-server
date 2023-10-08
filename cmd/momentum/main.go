package main

import (
	"momentum-go-server/internal/store"
	"time"
)

func main() {
	store.ContentDB()
	store.CreateUser()
	// store.GetUser()
	// routes.SetupRoutes()
	time.Sleep(10 * time.Second)
}
