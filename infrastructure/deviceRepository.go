package infrastructure

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"os"
	"theorem-prototype/domain"
	"time"
)

type deviceRepositoryImplementation struct {
	sqlConnection *sql.DB
}

func NewDeviceRepository() (domain.DeviceRepository, error) {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, username, password, dbName)

	log.Printf("Opening connection to database with parameters: %s", dbInfo)

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

func (repository *deviceRepositoryImplementation) GetList() ([]*domain.Device, error) {
	var (
		serialNumber     string
		registrationDate time.Time
		firmwareVersion  string
	)
	devices := make([]*domain.Device, 0)

	query := `SELECT "SerialNumber", "RegistrationDate", "FirmwareVersion" FROM device."Devices"`
	rows, err := repository.sqlConnection.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting list of devices")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&serialNumber, &registrationDate, &firmwareVersion)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting list of devices")
		}

		device := domain.Device{SerialNumber: serialNumber, RegistrationDate: registrationDate, FirmwareVersion: firmwareVersion}
		devices = append(devices, &device)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting list of devices")
	}

	return devices, nil
}
