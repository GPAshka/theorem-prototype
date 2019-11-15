package infrastructure

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"theorem-prototype/config"
	"theorem-prototype/domain"
	"time"
)

type deviceRepositoryImplementation struct {
	sqlConnection *sql.DB
}

func NewDeviceRepository() (domain.DeviceRepository, error) {
	dbInfo := config.GetDataBaseInfo()

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, errors.Wrap(err, "error while opening database connection")
	}

	return &deviceRepositoryImplementation{sqlConnection: db}, nil
}

func (repository *deviceRepositoryImplementation) Add(device domain.Device) error {
	query := `INSERT INTO device."Devices" ("SerialNumber", "RegistrationDate", "FirmwareVersion") VALUES($1, $2, $3)`

	_, err := repository.sqlConnection.Exec(query, device.SerialNumber, device.RegistrationDate, device.FirmwareVersion)
	if err != nil {
		return errors.Wrap(err, "error while adding device to database")
	}

	return nil
}

func (repository *deviceRepositoryImplementation) Get(serialNumber string) (*domain.Device, error) {
	var device domain.Device

	query := `SELECT "SerialNumber", "RegistrationDate", "FirmwareVersion" FROM device."Devices" WHERE "SerialNumber" = $1`
	err := repository.sqlConnection.QueryRow(query, serialNumber).Scan(&device.SerialNumber, &device.RegistrationDate,
		&device.FirmwareVersion)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "error while getting device by serial number")
		}
	}

	return &device, nil
}

func (repository *deviceRepositoryImplementation) GetList(serialNumber string) ([]*domain.Device, error) {
	var (
		srNumber         string
		registrationDate time.Time
		firmwareVersion  string
	)
	devices := make([]*domain.Device, 0)

	query := `SELECT "SerialNumber", "RegistrationDate", "FirmwareVersion" FROM device."Devices" WHERE "SerialNumber" = $1 OR $1 = ''`
	rows, err := repository.sqlConnection.Query(query, serialNumber)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting list of devices")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&srNumber, &registrationDate, &firmwareVersion)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting list of devices")
		}

		device := domain.Device{SerialNumber: srNumber, RegistrationDate: registrationDate, FirmwareVersion: firmwareVersion}
		devices = append(devices, &device)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting list of devices")
	}

	return devices, nil
}
