package persistence

import (
	"time"

	"gorm.io/gorm"
)

// BaseGormModel contains common gorm model fields
type BaseGormModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate do field initialization before insert a row.
func (base *BaseGormModel) BeforeCreate(tx *gorm.DB) error {
	currentDateTime, err := currentDateTime()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("CreatedAt", currentDateTime)
	tx.Statement.SetColumn("UpdatedAt", currentDateTime)
	return nil
}

func currentDateTime() (time.Time, error) {
	return time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
