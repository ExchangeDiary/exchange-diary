package repository

import (
	"github.com/exchange-diary/domain/entity"
)

type RoomRepository interface {
	Create(room *entity.Room) (*entity.Room, error)
	GetByID(id int) (*entity.Room, error)
	GetAll(offset, limit int) (*entity.Rooms, error)
	GetAllByAccountId(accountId, offset, limit int) (*entity.Rooms, error)
	Update(room *entity.Room) (*entity.Room, error)
	Delete(room *entity.Room) error
}
