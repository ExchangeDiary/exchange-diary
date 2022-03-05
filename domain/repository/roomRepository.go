package repository

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
)

// RoomRepository ...
type RoomRepository interface {
	Create(room *entity.Room) (*entity.Room, error)
	GetByID(id uint) (*entity.Room, error)
	GetAll(accountID uint, roomIDs []uint, limit, offset uint) (*entity.Rooms, error)
	Update(room *entity.Room) (*entity.Room, error)
	Delete(room *entity.Room) error
}
