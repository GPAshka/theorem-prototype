package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"theorem-prototype/domain"
	"theorem-prototype/infrastructure"
)

type service struct {
	repository domain.DeviceRepository
	router     *mux.Router
}

func newService() (*service, error) {
	rep, err := infrastructure.NewDeviceRepository()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	serv := &service{router: router, repository: rep}
	return serv, nil
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
	serv, err := newService()
	if err != nil {
		log.Fatal(err)
		return
	}

	serv.router.HandleFunc("/api/v1/hc", serv.handleHealthCheck()).Methods("GET")
	//serv.router.HandleFunc("/api/v1/device",).Methods("POST")

	port := os.Getenv("PORT")
	log.Printf("Starting service on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, serv.router))
}
