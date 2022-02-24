package service

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
)

// RoomService ...
type RoomService interface {
	Create(name, code, hint, theme string) (*entity.Room, error)
	Get(id uint) (*entity.Room, error)
	GetAllJoinedRooms(accountID, limit, offset uint) (*entity.Rooms, error)
	GetAll(limit, offset uint) (*entity.Rooms, error)
	Update(id uint, lastname, firstname string) (*entity.Room, error)
	Delete(id uint) error
}

type roomService struct {
	roomRepository repository.RoomRepository
}

// NewRoomService ...
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

func (rs *roomService) Get(id uint) (*entity.Room, error) {
	room, err := rs.roomRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *roomService) GetAllJoinedRooms(accountID, limit, offset uint) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}
func (rs *roomService) GetAll(limit, offset uint) (*entity.Rooms, error) {
	rooms, err := rs.roomRepository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
func (rs *roomService) Update(id uint, lastname, firstname string) (*entity.Room, error) {
	return &entity.Room{}, nil
}
func (rs *roomService) Delete(id uint) error {
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
