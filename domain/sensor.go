package domain

import (
	"github.com/pkg/errors"
	"time"
)

type SensorRepository interface {
	AddSensorData(data *SensorData) error
	AddBulkSensorData(data []*SensorData) error
	GetSensorData(serialNumber string, date time.Time) ([]*SensorData, error)
}

type SensorData struct {
	DeviceSerialNumber string
	Date               time.Time
	Temperature        float32
	AirHumidity        float32
	CarbonMonoxide     float32
	HealthStatus       string
}

func (sensor *SensorData) Validate() error {
	if sensor.DeviceSerialNumber == "" {
		return errors.New("device serial number is required")
	}

	if sensor.Date.IsZero() {
		return errors.New("device sensor adding date is required")
	}

	if sensor.HealthStatus == "" {
		return errors.New("device sensor health status is required")
	}

	if len(sensor.HealthStatus) > 150 {
		return errors.New("device sensor health status should be less than 150 characters")
	}

	return nil
}
