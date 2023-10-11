package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"momentum-go-server/internal/handlers"
	"momentum-go-server/internal/middlware"
)

type APIServer struct {
	listenPort string
}

func NewAPIServer(listenPort string) *APIServer {
	return &APIServer{
		listenPort: listenPort,
	}
}

func (s *APIServer) Run() {

	r := mux.NewRouter()

	r.HandleFunc("/auth/signup", handlers.SignUpHandler).Methods("POST")
	r.HandleFunc("/auth/signin", handlers.SignInHandler).Methods("POST")
	r.HandleFunc("/", middlware.VerifyJWT(handlers.Home)).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server is listening on port: ", s.listenPort)
	http.ListenAndServe(s.listenPort, nil)
}
