package routes

import (
	"fmt"
	"net/http"

	h "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"momentum-go-server/internal/handlers"
	"momentum-go-server/internal/middlware"
	"momentum-go-server/internal/utils"
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
	r.HandleFunc("/setting/{type}", middlware.VerifyJWT(handlers.UpdateSettings)).Methods("PUT")

	headersOk := h.AllowedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := h.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := h.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT"})
	exposedHeadersOk := h.ExposedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})

	http.Handle("/", r)
	fmt.Println("Server is listening on port: ", s.listenPort)
	utils.InfoLogger.Println("Server is listening on port: ", s.listenPort)
	http.ListenAndServe(s.listenPort, h.CORS(originsOk, headersOk, methodsOk, exposedHeadersOk)(r))
}
