package domain

import (
	"errors"
	"time"
)

type DeviceRepository interface {
	Get(serialNumber string) (*Device, error)
	GetList(serialNumber string) ([]*Device, error)
	Add(device Device) error
}

type Device struct {
	SerialNumber     string
	RegistrationDate time.Time
	FirmwareVersion  string
}

func (device *Device) Validate() error {
	if device.SerialNumber == "" {
		return errors.New("device serial number is required")
	}

	if device.FirmwareVersion == "" {
		return errors.New("device firmware version is required")
	}

	if device.RegistrationDate.IsZero() {
		return errors.New("device registration date is required")
	}

	return nil
}
