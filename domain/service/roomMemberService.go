package service

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
)

// RoomMemberService ...
type RoomMemberService interface {
	Add(roomID, accountID uint) (*entity.RoomMember, error)
	Get(roomID, accountID uint) (*entity.RoomMember, error)
	GetAll() (*entity.RoomMembers, error)
	Delete(roomID, accountID uint) error
}

type roomMemberService struct {
	roomMemberRepository repository.RoomMemberRepository
}

// NewRoomMemberService ...
func NewRoomMemberService(r repository.RoomMemberRepository) RoomMemberService {
	return &roomMemberService{roomMemberRepository: r}
}

func (rms *roomMemberService) Add(roomID, accountID uint) (*entity.RoomMember, error) {
	roomMember, err := entity.NewRoomMember(roomID, accountID)
	if err != nil {
		return nil, err
	}
	createdRoomMember, err := rms.roomMemberRepository.Create(roomMember)
	if err != nil {
		return nil, err
	}
	return createdRoomMember, nil
}

func (rms *roomMemberService) Get(roomID, accountID uint) (*entity.RoomMember, error) {
	roomMember, err := rms.roomMemberRepository.GetByUnq(roomID, accountID)
	if err != nil {
		return nil, err
	}
	return roomMember, nil
}

func (rms *roomMemberService) GetAll() (*entity.RoomMembers, error) {
	roomMembers, err := rms.roomMemberRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return roomMembers, nil
}

func (rms *roomMemberService) Delete(roomID, accountID uint) error {
	roomMember, err := rms.Get(roomID, accountID)
	if err != nil {
		return err
	}
	err = rms.roomMemberRepository.Delete(roomMember)
	if err != nil {
		return err
	}
	return nil
}
