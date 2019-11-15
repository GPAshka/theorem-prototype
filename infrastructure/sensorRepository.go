package infrastructure

import (
	"database/sql"
	"github.com/pkg/errors"
	"theorem-prototype/config"
	"theorem-prototype/domain"
	"time"
)

type sensorRepositoryImplementation struct {
	sqlConnection *sql.DB
}

func NewSensorRepository() (domain.SensorRepository, error) {
	dbInfo := config.GetDataBaseInfo()

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, errors.Wrap(err, "error while opening database connection")
	}

	return &sensorRepositoryImplementation{sqlConnection: db}, nil
}

func (repository *sensorRepositoryImplementation) AddSensorData(data *domain.SensorData) error {
	query := `INSERT INTO device."SensorData" ("DeviceSerialNumber", "Date", "Temperature", "AirHumidity", "CarbonMonoxide", "HealthStatus") 
				VALUES($1, $2, $3, $4, $5, $6)`

	_, err := repository.sqlConnection.Exec(query, data.DeviceSerialNumber, data.Date, data.Temperature, data.AirHumidity,
		data.CarbonMonoxide, data.HealthStatus)
	if err != nil {
		return errors.Wrap(err, "error while adding device sensor data to database")
	}

	return nil
}

func (repository *sensorRepositoryImplementation) GetSensorData(serialNumber string) ([]*domain.SensorData, error) {
	var (
		deviceSerialNumber string
		date               time.Time
		temperature        float32
		airHumidity        float32
		carbonMonoxide     float32
		healthStatus       string
	)

	sensorData := make([]*domain.SensorData, 0)

	query := `SELECT "Date", "Temperature", "AirHumidity", "CarbonMonoxide", "HealthStatus", "DeviceSerialNumber"
				FROM device."SensorData"
				WHERE "DeviceSerialNumber" = $1;`

	rows, err := repository.sqlConnection.Query(query, serialNumber)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting sensor data for device")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&date, &temperature, &airHumidity, &carbonMonoxide, &healthStatus, &deviceSerialNumber)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting sensor data for device")
		}

		data := domain.SensorData{
			DeviceSerialNumber: deviceSerialNumber,
			Date:               date,
			Temperature:        temperature,
			AirHumidity:        airHumidity,
			CarbonMonoxide:     carbonMonoxide,
			HealthStatus:       healthStatus,
		}
		sensorData = append(sensorData, &data)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting sensor data for device")
	}

	return sensorData, nil
}
