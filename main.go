package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type service struct {
	router *mux.Router
}

func newService() *service {
	router := mux.NewRouter()
	serv := &service{router: router}
	return serv
}

func (s *service) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	serv := newService()

	serv.router.HandleFunc("/api/v1/hc", serv.handleHealthCheck()).Methods("GET")

	port := os.Getenv("PORT")
	log.Printf("Starting service on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, serv.router))
}
