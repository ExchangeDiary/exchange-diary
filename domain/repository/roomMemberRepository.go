package repository

import (
	"github.com/exchange-diary/domain/entity"
)

// RoomMemberRepository ...
type RoomMemberRepository interface {
	Create(roomMember *entity.RoomMember) (*entity.RoomMember, error)
	GetByUnq(roomID, accountID uint) (*entity.RoomMember, error)
	GetAll() (*entity.RoomMembers, error)
	Delete(roomMember *entity.RoomMember) error
}
