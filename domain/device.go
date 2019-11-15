package domain

import "time"

type AirCondDeviceRepository interface {
	Add(device AirCondDevice) error
}

type AirCondDevice struct {
	Id               int
	SerialNumber     string
	RegistrationDate time.Time
	FirmwareVersion  string
}
