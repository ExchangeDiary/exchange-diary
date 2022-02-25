package service

import (
	"errors"

	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
)

// RoomService ...
type RoomService interface {
	Create(masterID uint, name, code, hint, theme string, period uint8) (*entity.Room, error)
	Get(id uint) (*entity.Room, error)
	GetAllJoinedRooms(accountID, limit, offset uint) (*entity.Rooms, error)
	GetAll(limit, offset uint) (*entity.Rooms, error)
	Update(id uint, lastname, firstname string) (*entity.Room, error)
	Delete(id uint) error
	VerifyCode(id uint, code string) (bool, error)
	JoinRoom(roomID, accountID uint) error
}

type roomService struct {
	roomRepository repository.RoomRepository
}

// NewRoomService ...
func NewRoomService(roomRepository repository.RoomRepository) RoomService {
	return &roomService{roomRepository: roomRepository}
}

func (rs *roomService) Create(masterID uint, name, code, hint, theme string, period uint8) (*entity.Room, error) {
	room, err := entity.NewRoom(masterID, name, code, hint, theme, period)
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

// Room Master + RoomMember table
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

func (rs *roomService) VerifyCode(id uint, code string) (bool, error) {
	room, err := rs.Get(id)
	if err != nil {
		return false, err
	}
	if room.Code != code {
		return false, errors.New("Invalid code is given")
	}
	return true, nil
}

// update room.Orders
// add roomMember row
func (rs *roomService) JoinRoom(roomID, accountID uint) error {
	return nil
}
