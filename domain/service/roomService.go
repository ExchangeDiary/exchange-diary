package service

import (
	"errors"
	"fmt"

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
	JoinRoom(id, accountID uint, code string) (bool, error)
	LeaveRoom(id, accountID uint) error
}

type roomService struct {
	roomRepository    repository.RoomRepository
	roomMemberService RoomMemberService
}

// NewRoomService ...
func NewRoomService(rr repository.RoomRepository, rms RoomMemberService) RoomService {
	return &roomService{
		roomRepository:    rr,
		roomMemberService: rms,
	}
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

func (rs *roomService) JoinRoom(id, accountID uint, code string) (bool, error) {
	// get a room
	room, err := rs.Get(id)
	if err != nil {
		return false, err
	}
	if room.IsAlreadyJoined(accountID) {
		return false, errors.New("Already joined room")
	}
	// validate code
	if room.Code != code {
		return false, errors.New("Invalid code is given")
	}

	// ===TODO: tx start===
	// 	1. add roomMember
	if _, err := rs.roomMemberService.Add(id, accountID); err != nil {
		return false, err
	}
	// 	2. append room.Orders
	// TODO: Update JSON list field
	room.AppendMember(accountID)
	if _, err := rs.roomRepository.Update(room); err != nil {
		return false, err
	}
	// ===TODO: tx end===

	return true, nil
}

func (rs *roomService) LeaveRoom(id, accountID uint) error {
	// get a room
	room, err := rs.Get(id)
	if err != nil {
		return err
	}
	if !room.IsAlreadyJoined(accountID) {
		return nil
	}
	if room.IsMaster(accountID) {
		return rs.doMasterLeaveProcess(room, accountID)
	}
	return rs.doMemberLeaveProcess(room, accountID)
}

func (rs *roomService) doMasterLeaveProcess(room *entity.Room, accountID uint) error {
	fmt.Println("doMasterLeaveProcess")
	// 다이어리방에 한명만 존재할 경우: 다이어리방 제거
	if len(room.Orders) == 1 {
		if err := rs.roomRepository.Delete(room); err != nil {
			return err
		}
	}
	// 새로운 마스터 선출
	if err := room.ChangeMaster(); err != nil {
		return err
	}
	// TODO: TurnAccountID 변경
	// Update room
	_, err := rs.roomRepository.Update(room)
	return err
}

func (rs *roomService) doMemberLeaveProcess(room *entity.Room, accountID uint) error {
	fmt.Println("doMemberLeaveProcess")
	// TODO: TurnAccountID 변경

	// room.Order에서 빼기
	if _, err := room.RemoveMember(accountID); err != nil {
		return err
	}

	// Update room
	if _, err := rs.roomRepository.Update(room); err != nil {
		return err
	}

	// roomMember에서 row 제거
	if err := rs.roomMemberService.Delete(room.ID, accountID); err != nil {
		return err
	}
	return nil
}
