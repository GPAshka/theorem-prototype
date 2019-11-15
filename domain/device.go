package domain

import "time"

type DeviceRepository interface {
	Get(serialNumber string) (*Device, error)
	GetList() ([]*Device, error)
	Add(device Device) error
}

type Device struct {
	SerialNumber     string
	RegistrationDate time.Time
	FirmwareVersion  string
}
