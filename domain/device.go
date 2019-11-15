package domain

import "time"

type DeviceRepository interface {
	Add(device Device) error
}

type Device struct {
	SerialNumber     string
	RegistrationDate time.Time
	FirmwareVersion  string
}
