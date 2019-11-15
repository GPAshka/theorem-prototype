package infrastructure

import (
	"database/sql"
	"github.com/pkg/errors"
	"theorem-prototype/config"
	"theorem-prototype/domain"
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
