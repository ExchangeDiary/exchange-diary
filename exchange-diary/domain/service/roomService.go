package service

import (
	"github.com/ExchangeDiary_Server/exchange-diary/domain/entity"
	"github.com/ExchangeDiary_Server/exchange-diary/domain/repository"
)
type RoomService interface {
	Create(lastname, firstname string) (*entity.Room, error)
	Get(id int) (*entity.Room, error)
	GetAllJoinedRooms(accountId int) (*entity.Rooms, error)
	GetAll() (*entity.Rooms, error)
	Update(id int, lastname, firstname string) (*entity.Room, error)
	Delete(id int) error
}

type roomService struct {
	roomRepository repository.RoomRepository
}

func NewRoomService(roomRepository repository.RoomRepository) RoomService {
	return &roomService{roomRepository: roomRepository}	
}

func (roomService *roomService) Create(lastname, firstname string) (*entity.Room, error){ 
	return &entity.Room{}, nil
}
func (roomService *roomService) Get(id int) (*entity.Room, error){ 
	return &entity.Room{}, nil
}
func (roomService *roomService) GetAllJoinedRooms(accountId int) (*entity.Rooms, error){ 
	return &entity.Rooms{}, nil
}
func (roomService *roomService) GetAll() (*entity.Rooms, error){ 
	return &entity.Rooms{}, nil
}
func (roomService *roomService) Update(id int, lastname, firstname string) (*entity.Room, error){ 
	return &entity.Room{}, nil
}
func (roomService *roomService) Delete(id int) error{ 
	return nil
}