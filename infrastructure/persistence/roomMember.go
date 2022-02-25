package persistence

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// RoomMemberGorm is a db representation of entity.RoomMember
// "idx_room_account" is a unique key combined with (RoomID, AccountID)
type RoomMemberGorm struct {
	ID        uint        `gorm:"primaryKey"`
	RoomID    uint        `gorm:"column:room_id;uniqueIndex:idx_room_account"`
	Room      RoomGorm    `gorm:"foreignKey:RoomID"`
	AccountID uint        `gorm:"column:account_id;uniqueIndex:idx_room_account"`
	Account   AccountGorm `gorm:"foreignKey:AccountID"`
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

// GetByID ...
func (rmr *RoomMemberRepository) GetByID(id uint) (*entity.RoomMember, error) {
	dto := RoomMemberGorm{ID: id}
	if err := rmr.db.First(&dto).Error; err != nil {
		return nil, err
	}
	roomMember := new(entity.RoomMember)
	copier.Copy(&roomMember, &dto)
	return roomMember, nil
}

// GetAll ...
func (rmr *RoomMemberRepository) GetAll() (*entity.RoomMembers, error) {
	dto := RoomMemberGorms{}
	rmr.db.Find(&dto)
	roomMembers := new(entity.RoomMembers)
	copier.Copy(&roomMembers, &dto)
	return roomMembers, nil
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
