package repository

import (
	"github.com/exchange-diary/domain/entity"
)

// RoomRepository ...
type RoomRepository interface {
	Create(room *entity.Room) (*entity.Room, error)
	GetByID(id uint) (*entity.Room, error)
	GetAll(offset, limit uint) (*entity.Rooms, error)
	GetAllByAccountID(accountID, offset, limit uint) (*entity.Rooms, error)
	Update(room *entity.Room) (*entity.Room, error)
	Delete(room *entity.Room) error
}
