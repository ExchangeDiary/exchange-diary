package service

import (
	"github.com/ExchangeDiary_Server/exchange-diary/domain/entity"
	"github.com/ExchangeDiary_Server/exchange-diary/domain/repository"
)

type RoomService interface {
	Create(name, code, hint, theme string) (*entity.Room, error)
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

func (rs *roomService) Create(name, code, hint, theme string) (*entity.Room, error) {
	room, err := entity.NewRoom(name, code, hint, theme)
	if err != nil {
		return nil, err
	}
	createdRoom, err := rs.roomRepository.Create(room)
	if err != nil {
		return nil, err
	}
	return createdRoom, nil
}

func (rs *roomService) Get(id int) (*entity.Room, error) {
	return &entity.Room{}, nil
}
func (rs *roomService) GetAllJoinedRooms(accountId int) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}
func (rs *roomService) GetAll() (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}
func (rs *roomService) Update(id int, lastname, firstname string) (*entity.Room, error) {
	return &entity.Room{}, nil
}
func (rs *roomService) Delete(id int) error {
	return nil
}
