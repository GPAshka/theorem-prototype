package domain

import "time"

type SensorRepository interface {
	AddSensorData(data *SensorData) error
}

type SensorData struct {
	DeviceSerialNumber string
	Date               time.Time
	Temperature        float32
	AirHumidity        float32
	CarbonMonoxide     float32
	HealthStatus       string
}
