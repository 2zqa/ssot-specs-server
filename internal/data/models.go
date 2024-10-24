package data

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// List of all models
type Models struct {
	Devices DeviceModel
	Search  SearchModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Devices: DeviceModel{DB: db},
		Search:  SearchModel{DB: db},
	}
}

func convertSQLErrors(err error) error {
	switch {
	case err == nil:
		return nil
	case err.Error() == "ERROR: duplicate key value violates unique constraint \"devices_pkey\" (SQLSTATE 23505)":
		return ErrDuplicateDeviceID
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrRecordNotFound
	default:
		return err
	}
}

func WithPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Specs.Motherboard").
		Preload("Specs.CPU").
		Preload("Specs.Disks.Partitions").
		Preload("Specs.Network.Interfaces.Driver").
		Preload("Specs.Bios").
		Preload("Specs.Memory.SwapDevices").
		Preload("Specs.Kernel").
		Preload("Specs.Release").
		Preload("Specs.DIMMs").
		Preload("Specs.Virtualization").
		Preload("Specs.OEM")
}
