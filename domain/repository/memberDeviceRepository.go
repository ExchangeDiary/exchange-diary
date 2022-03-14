package repository

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
)

// MemberDeviceRepository ...
type MemberDeviceRepository interface {
	CreateIfNotExist(memberID uint, token string) (memberDevice *entity.MemberDevice, err error)
	Get(token string) (memberDevice *entity.MemberDevice, err error)
	GetAllTokens(memberID uint) (tokens []string)
	Delete(memberDevice *entity.MemberDevice) error
}
