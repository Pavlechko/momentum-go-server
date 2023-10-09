package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"momentum-go-server/internal"
	"momentum-go-server/internal/handlers"
)

func SetupRoutes() {
	m := internal.Message()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the home page!", m)
	})

	r.HandleFunc("/auth/signup", handlers.SignUpHandler).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
