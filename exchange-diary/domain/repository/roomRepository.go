package repository
import (
	"github.com/ExchangeDiary_Server/exchange-diary/domain/entity"
)

type RoomRepository interface {
	Create(room *entity.Room) (*entity.Room, error)
	GetByID(id int) (*entity.Room, error)
	GetAll() (*entity.Rooms, error)
	GetAllByAccountId(accountId int) (*entity.Rooms, error)
	Update(room *entity.Room) (*entity.Room, error)
	Delete(room *entity.Room) error
}