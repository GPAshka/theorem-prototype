package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"theorem-prototype/app"
	"theorem-prototype/domain"
	"theorem-prototype/infrastructure"
	"theorem-prototype/utils"
	"time"
)

type service struct {
	deviceRepository domain.DeviceRepository
	sensorRepository domain.SensorRepository
	router           *mux.Router
}

func newService() (*service, error) {
	deviceRep, err := infrastructure.NewDeviceRepository()
	if err != nil {
		return nil, err
	}

	sensorRep, err := infrastructure.NewSensorRepository()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	serv := &service{router: router, deviceRepository: deviceRep, sensorRepository: sensorRep}
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

		//validate request
		err = device.Validate()
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		//check if device with specified serial number already exists
		existingDevice, err := s.deviceRepository.Get(device.SerialNumber)
		if err != nil || existingDevice != nil {
			if existingDevice != nil {
				err = errors.New(fmt.Sprintf("Device with serial number '%s' already registered", device.SerialNumber))
			}

			utils.RespondError(w, err)
			return
		}

		//add new device
		err = s.deviceRepository.Add(device)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, nil)
	}
}

func (s *service) handleAddSensorData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//decode HTTP request body
		var sensorData domain.SensorData
		err := utils.DecodeRequest(r.Body, &sensorData)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		//check if device with specified serial number exists
		existingDevice, err := s.deviceRepository.Get(sensorData.DeviceSerialNumber)
		if err != nil || existingDevice == nil {
			if existingDevice == nil {
				err = errors.New(fmt.Sprintf("Device with serial number '%s' is not registered", sensorData.DeviceSerialNumber))
			}

			utils.RespondError(w, err)
			return
		}

		//add sensor data for device
		err = s.sensorRepository.AddSensorData(&sensorData)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, nil)
	}
}

func (s *service) handleAddBulkSensorData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//decode HTTP request body
		var sensorData []*domain.SensorData
		err := utils.DecodeRequest(r.Body, &sensorData)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		//add bulk sensor data
		err = s.sensorRepository.AddBulkSensorData(sensorData)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, nil)
	}
}

func (s *service) handleGetDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serialNumber := r.FormValue("serialNumber")

		devices, err := s.deviceRepository.GetList(serialNumber)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, devices)
	}
}

func (s *service) handleGetSensorData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//read request parameters
		vars := mux.Vars(r)
		serialNumber := vars["serialNumber"]
		dateParam := r.FormValue("date")

		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			utils.RespondError(w, errors.Wrap(err, "error while parsing 'date' query parameter"))
			return
		}

		sensorData, err := s.sensorRepository.GetSensorData(serialNumber, date)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		utils.RespondSuccess(w, sensorData)
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
	serv.router.HandleFunc("/api/v1/devices/sensors", serv.handleAddSensorData()).Methods("POST")
	serv.router.HandleFunc("/api/v1/devices/sensors/bulk", serv.handleAddBulkSensorData()).Methods("POST")

	serv.router.HandleFunc("/api/v1/devices", serv.handleGetDevices()).Methods("GET")
	serv.router.HandleFunc("/api/v1/devices/{serialNumber}/sensors", serv.handleGetSensorData()).Methods("GET")

	serv.router.Use(app.JwtAuthenticationMiddleware)

	port := os.Getenv("PORT")
	log.Printf("Starting service on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, serv.router))
}
