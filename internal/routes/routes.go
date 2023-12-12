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
	handler    handlers.IHandler
}

func NewAPIServer(listenPort string, handler handlers.IHandler) *APIServer {
	return &APIServer{
		listenPort: listenPort,
		handler:    handler,
	}
}

func (s *APIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/signup", s.handler.SignUpHandler).Methods("POST")
	r.HandleFunc("/auth/signin", s.handler.SignInHandler).Methods("POST")
	r.HandleFunc("/", middlware.VerifyJWT(s.handler.Home)).Methods("GET")
	r.HandleFunc("/setting/{type}", middlware.VerifyJWT(s.handler.UpdateSettings)).Methods("PUT")

	headersOk := h.AllowedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := h.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := h.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT"})
	exposedHeadersOk := h.ExposedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})

	http.Handle("/", r)
	fmt.Println("Server is listening on port: ", s.listenPort)
	utils.InfoLogger.Println("Server is listening on port: ", s.listenPort)
	server := &http.Server{
		Addr:    s.listenPort,
		Handler: h.CORS(originsOk, headersOk, methodsOk, exposedHeadersOk)(r),
	}
	err := server.ListenAndServe()
	if err != nil {
		utils.ErrorLogger.Println("Error listening server:", err.Error())
	}
}
