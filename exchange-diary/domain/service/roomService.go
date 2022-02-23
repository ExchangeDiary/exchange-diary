package service

import (
	"github.com/ExchangeDiary_Server/exchange-diary/domain/entity"
	"github.com/ExchangeDiary_Server/exchange-diary/domain/repository"
)

type RoomService interface {
	Create(name, code, hint, theme string) (*entity.Room, error)
	Get(id int) (*entity.Room, error)
	GetAllJoinedRooms(accountId, limit, offset int) (*entity.Rooms, error)
	GetAll(limit, offset int) (*entity.Rooms, error)
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
	room, err := rs.roomRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *roomService) GetAllJoinedRooms(accountId, limit, offset int) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}
func (rs *roomService) GetAll(limit, offset int) (*entity.Rooms, error) {
	rooms, err := rs.roomRepository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
func (rs *roomService) Update(id int, lastname, firstname string) (*entity.Room, error) {
	return &entity.Room{}, nil
}
func (rs *roomService) Delete(id int) error {
	room, err := rs.Get(id)
	if err != nil {
		return err
	}

	err = rs.roomRepository.Delete(room)
	if err != nil {
		return err
	}
	return nil
}
