package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
)

// RoomMemberService ...
type RoomMemberService interface {
	Add(roomID, accountID uint) (*entity.RoomMember, error)
	Get(roomID, accountID uint) (*entity.RoomMember, error)
	GetAllMemberRoomIDs(accountID uint) ([]uint, error)
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

func (rms *roomMemberService) GetAllMemberRoomIDs(accountID uint) (roomIDs []uint, err error) {
	roomIDs, err = rms.roomMemberRepository.GetAllRoomIDsByMemberID(accountID)
	if err != nil {
		return nil, err
	}
	return roomIDs, nil
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
