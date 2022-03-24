package repository

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
)

// AlarmRepository ...
type AlarmRepository interface {
	Create(alarm *entity.Alarm) (*entity.Alarm, error)
	GetAll(accountID uint) (*entity.Alarms, error)
}
