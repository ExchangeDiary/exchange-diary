package persistence

import (
	"fmt"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// MemberDeviceGorm ...
type MemberDeviceGorm struct {
	ID       uint       `gorm:"primaryKey"`
	MemberID uint       `gorm:"column:member_id"`
	Member   MemberGorm `gorm:"index;column:member_id;constraint:OnDelete:CASCADE;"`
	// varchar(512)사용 이유: https://stackoverflow.com/questions/11668761/gcm-max-length-for-registration-id
	DeviceToken string `gorm:"column:device_token;type:varchar(512);uniqueIndex:idx_device_token;not null"`
	BaseGormModel
}

// TableName define gorm table name
func (MemberDeviceGorm) TableName() string {
	return "member_devices"
}

// MemberDeviceGorms define list of MemberDeviceGorm
type MemberDeviceGorms []MemberDeviceGorm

// MemberDeviceRepository ...
type MemberDeviceRepository struct {
	db *gorm.DB
}

// NewMemberDeviceRepository ...
func NewMemberDeviceRepository(db *gorm.DB) repository.MemberDeviceRepository {
	return &MemberDeviceRepository{db: db}
}

// ToEntity ...
func (dto *MemberDeviceGorm) ToEntity() *entity.MemberDevice {
	memberDevice := new(entity.MemberDevice)
	copier.Copy(&memberDevice, &dto)
	return memberDevice
}

// CreateIfNotExist ...
func (mdr *MemberDeviceRepository) CreateIfNotExist(memberID uint, token string) (memberDevice *entity.MemberDevice, err error) {
	if memberDevice, err = mdr.Get(token); err != nil {
		return mdr.create(memberID, token)
	}
	return memberDevice, nil
}

func (mdr *MemberDeviceRepository) create(memberID uint, token string) (memberDevice *entity.MemberDevice, err error) {
	dto := MemberDeviceGorm{MemberID: memberID, DeviceToken: token}
	if err := mdr.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	return dto.ToEntity(), nil
}

// Get ...
func (mdr *MemberDeviceRepository) Get(token string) (memberDevice *entity.MemberDevice, err error) {
	dto := MemberDeviceGorm{}
	if err := mdr.db.Where("device_token = ?", token).First(&dto).Error; err != nil {
		return nil, err
	}
	return dto.ToEntity(), nil
}

// GetAllTokens ...
func (mdr *MemberDeviceRepository) GetAllTokens(memberID uint) (tokens []string, err error) {
	dto := MemberDeviceGorms{}
	mdr.db.Select("DeviceToken").Where("member_id = ?", memberID).Find(&dto)
	if len(dto) < 1 {
		return nil, fmt.Errorf(fmt.Sprintf("There is no device tokens. memberID: %d", memberID))
	}
	for _, memberDeviceGorm := range dto {
		tokens = append(tokens, memberDeviceGorm.DeviceToken)
	}
	return tokens, nil
}

// GetAllMemberTokens ...
func (mdr *MemberDeviceRepository) GetAllMemberTokens(memberIDs []uint) (tokens []string, err error) {
	dto := MemberDeviceGorms{}
	mdr.db.Select("DeviceToken").Where("member_id IN (?)", memberIDs).Find(&dto)
	if len(dto) < 1 {
		return nil, fmt.Errorf(fmt.Sprintf("There is no device tokens. memberIDs: %d", memberIDs))
	}

	for _, memberDeviceGorm := range dto {
		tokens = append(tokens, memberDeviceGorm.DeviceToken)
	}
	return tokens, nil
}

// Delete ...
func (mdr *MemberDeviceRepository) Delete(memberDevice *entity.MemberDevice) error {
	dto := MemberDeviceGorm{}
	copier.Copy(&dto, &memberDevice)
	if err := mdr.db.Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBatch ...
func (mdr *MemberDeviceRepository) DeleteBatch(tokens []string) error {
	dto := MemberDeviceGorm{}
	if err := mdr.db.Where("device_token IN (?)", tokens).Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}
