package data

import (
	"context"
	"errors"
	"time"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const databaseTimeout time.Duration = 3 * time.Second

var (
	ErrDuplicateDeviceID = errors.New("duplicate device")
)

type Device struct {
	ID        uuid.UUID `json:"uuid" gorm:"type:uuid"`
	UpdatedAt time.Time `json:"updated_at"`
	Specs     Specs     `json:"specs,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

func (d Device) String() string {
	return d.ID.String() + delimiter + d.Specs.String()
}

func ValidateDevice(v *validator.Validator, device *Device) {
	v.Check(device.ID != uuid.Nil, "uuid", "must be provided")
	ValidateSpecs(v, &device.Specs)
}

type DeviceModel struct {
	DB *gorm.DB
}

func (m DeviceModel) Insert(device *Device, s SearchModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()
	err := convertSQLErrors(m.DB.WithContext(ctx).Create(device).Error)
	searchData := device.String()
	s.Insert(&Search{device.ID, searchData})
	return err
}

func (m DeviceModel) Update(device *Device, s SearchModel) error {
	// https://gorm.io/docs/transactions.html#Transaction
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		err := m.Delete(device.ID, s)
		if err != nil {
			return err
		}
		return m.Insert(device, s)
	})

	return convertSQLErrors(err)
}

func (m DeviceModel) Get(uuid uuid.UUID) (*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()
	device := &Device{
		ID: uuid,
	}
	err := convertSQLErrors(m.DB.Scopes(WithPreloads).WithContext(ctx).Take(device).Error)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (m DeviceModel) Delete(uuid uuid.UUID, s SearchModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()
	s.Delete(uuid)
	return convertSQLErrors(m.DB.WithContext(ctx).Delete(&Device{}, uuid).Error)
}

func (m DeviceModel) GetAll(searchTerms string, filters Filters) ([]*Device, Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()
	devices := []*Device{}
	metadata := Metadata{}
	err := convertSQLErrors(m.DB.
		Scopes(
			WithPreloads,
			filters.search(searchTerms),
			filters.sort(),
			filters.paginate(m.DB, searchTerms, &metadata),
		).
		WithContext(ctx).
		Find(&devices).
		Error)
	if err != nil {
		return nil, metadata, err
	}
	return devices, metadata, nil
}
