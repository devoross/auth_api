package server

import (
	"auth_api/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Router *mux.Router
	Port   string
}

func New(p string) *Server {
	log.Println("msg=\"setting up web server...\", app=\"auth_api\", level=\"info\"")
	s := &Server{
		Router: mux.NewRouter(),
		Port:   p,
	}
	u := users.NewRedis()

	s.Router.Use(users.PrometheusMiddleware)

	s.Router.HandleFunc("/health", apiStatus)
	s.Router.HandleFunc("/api/auth/login", u.LoginHandler)
	s.Router.HandleFunc("/api/auth/register", u.SignupHandler)

	s.Router.Handle("/metrics", promhttp.Handler())

	log.Println("msg=\"server setup complete...\", app=\"auth_api\", level=\"info\"")
	return s
}

func apiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"healthy": "yes"})
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.Router))
}
