package persistence

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// RoomMemberGorm is a db representation of entity.RoomMember
// "idx_room_account" is a unique key combined with (RoomID, AccountID)
type RoomMemberGorm struct {
	ID        uint       `gorm:"primaryKey"`
	RoomID    uint       `gorm:"column:room_id"`
	Room      RoomGorm   `gorm:"uniqueIndex:idx_room_account;column:room_id;constraint:OnDelete:CASCADE;"`
	AccountID uint       `gorm:"column:account_id"`
	Account   MemberGorm `gorm:"uniqueIndex:idx_room_account;column:account_id;constraint:OnDelete:CASCADE;"`
	BaseGormModel
}

// TableName define gorm table name
func (RoomMemberGorm) TableName() string {
	return "room_members"
}

// RoomMemberGorms define list of RoomMemberGorm
type RoomMemberGorms []RoomMemberGorm

// RoomMemberRepository is a impl of domain/repository/roomMemberRepository.go RoomMemberRepository interface
type RoomMemberRepository struct {
	db *gorm.DB
}

// NewRoomMemberRepository ...
func NewRoomMemberRepository(db *gorm.DB) repository.RoomMemberRepository {
	return &RoomMemberRepository{db: db}
}

// Create ...
func (rmr *RoomMemberRepository) Create(roomMember *entity.RoomMember) (*entity.RoomMember, error) {
	dto := RoomMemberGorm{}
	copier.Copy(&dto, &roomMember)
	if err := rmr.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	newRoomMember := new(entity.RoomMember)
	copier.Copy(&newRoomMember, &dto)
	return newRoomMember, nil
}

// GetByUnq func gets RoomMember row by unique_key(RoomID, AccountID)
func (rmr *RoomMemberRepository) GetByUnq(roomID, accountID uint) (*entity.RoomMember, error) {
	dto := RoomMemberGorm{RoomID: roomID, AccountID: accountID}
	if err := rmr.db.First(&dto).Error; err != nil {
		return nil, err
	}
	roomMember := new(entity.RoomMember)
	copier.Copy(&roomMember, &dto)
	return roomMember, nil
}

// GetAllRoomIDsByMemberID ...
func (rmr *RoomMemberRepository) GetAllRoomIDsByMemberID(memberID uint) (roomIDs []uint, err error) {
	dto := RoomMemberGorms{}
	rmr.db.Select("id").Where("account_id = ?", memberID).Find(&dto)
	for _, roomMemberGorm := range dto {
		roomIDs = append(roomIDs, roomMemberGorm.ID)
	}
	return roomIDs, err
}

// Delete ...
func (rmr *RoomMemberRepository) Delete(roomMember *entity.RoomMember) error {
	dto := RoomMemberGorm{}
	copier.Copy(&dto, &roomMember)
	if err := rmr.db.Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}
