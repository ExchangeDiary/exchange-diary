package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
)

// RoomMemberService ...
type RoomMemberService interface {
	Add(roomID, accountID uint) (*entity.RoomMember, error)
	Get(roomID, accountID uint) (*entity.RoomMember, error)
	GetAllRoomIDs(accountID uint) ([]uint, error)
	Delete(roomID, accountID uint) error
	PopulateRoomsMembers(rooms *entity.Rooms) (*entity.Rooms, error)
	PopulateRoomMembers(room *entity.Room) (*entity.Room, error)
}

type roomMemberService struct {
	roomMemberRepository repository.RoomMemberRepository
	memberRepository     repository.MemberRepository
}

// NewRoomMemberService ...
func NewRoomMemberService(r repository.RoomMemberRepository, mr repository.MemberRepository) RoomMemberService {
	return &roomMemberService{
		roomMemberRepository: r,
		memberRepository:     mr,
	}
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

func (rms *roomMemberService) GetAllRoomIDs(accountID uint) (roomIDs []uint, err error) {
	roomIDs, err = rms.roomMemberRepository.GetAllRoomIDs(accountID)
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

func (rms *roomMemberService) PopulateRoomsMembers(rooms *entity.Rooms) (*entity.Rooms, error) {
	populatedRooms := entity.Rooms{}
	// O(n)
	for _, room := range *rooms {
		populatedRoom, err := rms.PopulateRoomMembers(&room)
		if err != nil {
			return nil, err
		}
		populatedRooms = append(populatedRooms, *populatedRoom)
	}
	return &populatedRooms, nil
}

// roomMember.created_at 기준으로 populate
func (rms *roomMemberService) PopulateRoomMembers(room *entity.Room) (*entity.Room, error) {
	populatedRoom := entity.Room{}
	copier.Copy(&populatedRoom, &room)
	memberIDs, err := populatedRoom.MemberOnlyOrders()
	if err != nil {
		return nil, err
	}
	sortedMemberIDs, err := rms.roomMemberRepository.SortedMemberIDs(memberIDs)
	if err != nil {
		return nil, err
	}

	allIDs := append([]uint{room.MasterID}, sortedMemberIDs...)
	members, err := rms.memberRepository.GetAllByIDs(allIDs)
	if err != nil {
		return nil, err
	}
	populatedRoom.Members = members
	return &populatedRoom, nil
}

// TODO remove memberRepository dependency
// TODO: /v1/rooms/<:room_id>/orders 분리용
func (rms *roomMemberService) populateMembersByOrders(room *entity.Room) (*entity.Room, error) {
	if len(room.Orders) != 0 {
		// PopulateMembers
		members, err := rms.memberRepository.GetAllByIDs(room.Orders)
		if err != nil {
			return nil, err
		}
		(*room).Members = members
	}
	return room, nil
}
