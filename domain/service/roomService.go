package service

import (
	"fmt"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
)

// RoomService ...
type RoomService interface {
	Create(masterID uint, name, code, hint, theme string, period uint8) (*entity.Room, error)
	Get(id uint) (*entity.Room, error)
	GetAllJoinedRooms(accountID, limit, offset uint) (*entity.Rooms, error)
	Update(room *entity.Room) (*entity.Room, error)
	Delete(room *entity.Room) error
	JoinRoom(id, accountID uint, code string) (bool, error)
	LeaveRoom(id, accountID uint) error
}

type roomService struct {
	roomMemberService RoomMemberService
	roomRepository    repository.RoomRepository
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

func (rs *roomService) Get(id uint) (room *entity.Room, err error) {
	if room, err = rs.roomRepository.GetByID(id); err != nil {
		return nil, err
	}
	populatedRoom, err := rs.roomMemberService.PopulateRoomMembers(room)
	if err != nil {
		return nil, err
	}
	return populatedRoom, nil
}

// SELECT * FROM `rooms` WHERE id IN (memberRoomIDs) OR master_id = accountID ORDER BY  created_at desc  LIMIT limit OFFSET offset;
func (rs *roomService) GetAllJoinedRooms(accountID, limit, offset uint) (*entity.Rooms, error) {
	// O(1)
	memberRoomIDs, err := rs.roomMemberService.GetAllRoomIDs(accountID)
	if err != nil {
		return nil, err
	}
	// O(1)
	rooms, err := rs.roomRepository.GetAll(accountID, memberRoomIDs, limit, offset)
	if err != nil {
		return nil, err
	}

	populatedRooms, err := rs.roomMemberService.PopulateRoomsMembers(rooms)
	if err != nil {
		return nil, err
	}
	return populatedRooms, nil
}

func (rs *roomService) Update(room *entity.Room) (*entity.Room, error) {
	room, err := rs.roomRepository.Update(room)
	if err != nil {
		return nil, err
	}
	rs.roomMemberService.PopulateRoomMembers(room)
	return room, nil
}

func (rs *roomService) Delete(room *entity.Room) error {
	err := rs.roomRepository.Delete(room)
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
		return false, fmt.Errorf("Already joined room")
	}
	// validate code
	if room.Code != code {
		return false, fmt.Errorf("Invalid code is given")
	}

	// ===TODO: tx start===
	// 	1. add roomMember
	if _, err := rs.roomMemberService.Add(id, accountID); err != nil {
		return false, err
	}
	// 	2. append room.Orders
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
		return fmt.Errorf("cannot leave room because you are not a memeber of this room")
	}
	// TurnAccountID 변경
	if room.IsTurn(accountID) {
		room.NextTurn()
	}

	if room.IsMaster(accountID) {
		return rs.doMasterLeaveProcess(room, accountID)
	}
	return rs.doMemberLeaveProcess(room, accountID)
}

func (rs *roomService) doMasterLeaveProcess(room *entity.Room, accountID uint) error {
	// 다이어리방에 한명만 존재할 경우: 다이어리방 제거
	if len(room.Orders) == 1 {
		if err := rs.roomRepository.Delete(room); err != nil {
			return err
		}
		return nil
	}
	// 새로운 마스터 선출
	if err := room.ChangeMaster(); err != nil {
		return err
	}

	// 새로운 마스터를 멤버에서 제외
	if err := rs.roomMemberService.Delete(room.ID, room.MasterID); err != nil {
		return err
	}

	// Update room
	_, err := rs.roomRepository.Update(room)
	return err
}

func (rs *roomService) doMemberLeaveProcess(room *entity.Room, accountID uint) error {
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
