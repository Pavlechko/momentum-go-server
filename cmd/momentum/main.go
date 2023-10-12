package main

import (
	"momentum-go-server/internal/routes"
	"momentum-go-server/internal/store"
)

func main() {
	store.ContentDB()

	server := routes.NewAPIServer(":8080")
	server.Run()
	// time.Sleep(10 * time.Second)
}
