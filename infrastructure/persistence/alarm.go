package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// AlarmGorm is a db representation of entity.Member
type AlarmGorm struct {
	ID       uint       `gorm:"primaryKey"`
	MemberID uint       `gorm:"column:member_id"`
	Member   MemberGorm `gorm:"uniqueIndex:unq_room_member_code;column:member_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Code     string     `gorm:"uniqueIndex:unq_room_member_code;column:code;type:varchar(512)"`
	Title    string     `gorm:"column:title"`
	RoomID   uint       `gorm:"column:room_id"`
	Room     RoomGorm   `gorm:"uniqueIndex:unq_room_member_code;column:room_id;constraint:OnDelete:CASCADE;"`
	RoomName string     `gorm:"column:room_name"`
	Author   string     `gorm:"column:author"`
	AlarmAt  time.Time  `gorm:"column:alarm_at"`
}

// TableName define gorm table name
func (AlarmGorm) TableName() string {
	return "alarms"
}

// AlarmsGorm is a type that represents list of AlarmGorm
type AlarmsGorm []AlarmGorm

// AlarmRepository ...
type AlarmRepository struct {
	db *gorm.DB
}

// NewAlarmRepository ...
func NewAlarmRepository(db *gorm.DB) repository.AlarmRepository {
	return &AlarmRepository{db: db}
}

func (a *AlarmGorm) toEntity() *entity.Alarm {
	entity := new(entity.Alarm)
	copier.Copy(&entity, &a)
	return entity
}

func toDTO(a *entity.Alarm) *AlarmGorm {
	dto := new(AlarmGorm)
	copier.Copy(&dto, &a)
	return dto
}

// Create ...
func (ar *AlarmRepository) Create(alarm *entity.Alarm) (*entity.Alarm, error) {
	dto := toDTO(alarm)

	// is already exist
	if err := ar.deleteIfExist(alarm.UnqFields()); err != nil {
		return nil, err
	}

	if err := ar.db.Create(&dto).Error; err != nil {
		return nil, err
	}

	return dto.toEntity(), nil
}

func (ar *AlarmRepository) deleteIfExist(roomID, memberID uint, code vo.TaskCode) error {
	dto := AlarmGorm{}
	if err := ar.db.Where("room_id = ? AND member_id = ? AND code = ?", roomID, memberID, string(code)).First(&dto).Error; err != nil {
		fmt.Println("First에서 에러가 나긴함")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if err := ar.db.Delete(&dto).Error; err != nil {
		fmt.Println("Delete t에서 에러가 나긴함")
		return err
	}

	fmt.Println("Delete 정상적으로 처리됨")
	return nil
}

// GetAll ...
func (ar *AlarmRepository) GetAll(accountID uint) (*entity.Alarms, error) {
	dto := AlarmsGorm{}
	if err := ar.db.Where("member_id = ?", accountID).Order(" alarm_at desc ").Find(&dto).Error; err != nil {
		return nil, err
	}

	alarms := entity.Alarms{}
	for _, a := range dto {
		alarms = append(alarms, *a.toEntity())
	}
	return &alarms, nil
}
