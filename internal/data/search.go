package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Search struct {
	DeviceID uuid.UUID `json:"-" gorm:"type:uuid;primaryKey"`
	Data     string    `json:"-"`
}

type SearchModel struct {
	DB *gorm.DB
}

func (m SearchModel) Insert(search *Search) error {
	return convertSQLErrors(m.DB.Create(search).Error)
}

func (m SearchModel) Delete(uuid uuid.UUID) error {
	return convertSQLErrors(m.DB.Delete(&Search{}, uuid).Error)
}
