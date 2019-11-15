package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"theorem-prototype/domain"
	"theorem-prototype/infrastructure"
	"theorem-prototype/utils"
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

func (s *service) handleAddDevice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//decode HTTP request body
		var device domain.Device
		err := utils.DecodeRequest(r.Body, &device)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		//check if device with specified serial number already exists
		existingDevice, err := s.repository.Get(device.SerialNumber)
		if err != nil || existingDevice != nil {
			if existingDevice != nil {
				err = errors.New(fmt.Sprintf("Device with serial number '%s' already registered", device.SerialNumber))
			}

			utils.RespondError(w, err)
			return
		}

		//add new device
		err = s.repository.Add(device)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, nil)
	}
}

func (s *service) handleGetDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := s.repository.GetList()
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, devices)
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
	serv.router.HandleFunc("/api/v1/devices", serv.handleAddDevice()).Methods("POST")

	serv.router.HandleFunc("/api/v1/devices", serv.handleGetDevices()).Methods("GET")

	port := os.Getenv("PORT")
	log.Printf("Starting service on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, serv.router))
}
