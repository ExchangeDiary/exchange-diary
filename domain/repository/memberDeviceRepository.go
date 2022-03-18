package repository

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
)

// MemberDeviceRepository ...
type MemberDeviceRepository interface {
	CreateIfNotExist(memberID uint, token string) (memberDevice *entity.MemberDevice, err error)
	Get(token string) (memberDevice *entity.MemberDevice, err error)
	GetAllTokens(memberID uint) (tokens []string, err error)
	GetAllMemberTokens(memberIDs []uint) (tokens []string, err error)
	Delete(memberDevice *entity.MemberDevice) error
	DeleteBatch(tokens []string) error
}
