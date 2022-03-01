package repository

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
)

// RoomMemberRepository ...
type RoomMemberRepository interface {
	Create(roomMember *entity.RoomMember) (*entity.RoomMember, error)
	GetByUnq(roomID, accountID uint) (*entity.RoomMember, error)
	GetAllRoomIDsByMemberID(memberID uint) ([]uint, error)
	Delete(roomMember *entity.RoomMember) error
}
