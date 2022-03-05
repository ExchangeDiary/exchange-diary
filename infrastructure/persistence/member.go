package persistence

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// MemberGorm is a db representation of entity.Member
type MemberGorm struct {
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"column:email;uniqueIndex,not null"`
	Name              string `gorm:"column:name;not null"`
	ProfileURL        string `gorm:"column:profile_url"`
	AuthType          string `gorm:"column:auth_type"`
	TurnAlarmFlag     bool   `gorm:"column:turn_alarm_flag"`
	ActivityAlarmFlag bool   `gorm:"column:activity_alarm_flag"`
	BaseGormModel
}

// TableName define gorm table name
func (MemberGorm) TableName() string {
	return "members"
}

// MembersGorm is a type that represents list of MemberGorm
type MembersGorm []MemberGorm

// MemberRepository ...
type MemberRepository struct {
	db *gorm.DB
}

// NewMemberRepository ...
func NewMemberRepository(db *gorm.DB) repository.MemberRepository {
	return &MemberRepository{db: db}
}

// ToMemberEntity ...
func ToMemberEntity(memberDto *MemberGorm) *entity.Member {
	member := new(entity.Member)
	copier.Copy(&member, &memberDto)
	return member
}

// ToMemberDTO ...
func ToMemberDTO(member *entity.Member) *MemberGorm {
	memberDto := new(MemberGorm)
	copier.Copy(&memberDto, &member)
	return memberDto
}

// Create ...
func (r *MemberRepository) Create(member *entity.Member) (*entity.Member, error) {
	dto := ToMemberDTO(member)
	if err := r.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	return ToMemberEntity(dto), nil
}

// GetByEmail ...
func (r *MemberRepository) GetByEmail(email string) (*entity.Member, error) {
	dto := MemberGorm{}
	if err := r.db.Where("email = ?", email).Find(&dto).Error; err != nil {
		return nil, err
	}
	return ToMemberEntity(&dto), nil
}

// GetAllByIDs ...
func (r *MemberRepository) GetAllByIDs(ids []uint) (*entity.Members, error) {
	dto := MembersGorm{}
	if err := r.db.Where("id IN (?)", ids).Find(&dto).Error; err != nil {
		return nil, err
	}
	members := entity.Members{}
	for _, memberTO := range dto {
		members = append(members, *ToMemberEntity(&memberTO))
	}
	return &members, nil
}

// Update ...
func (r *MemberRepository) Update(member *entity.Member) (*entity.Member, error) {
	dto := ToMemberDTO(member)
	if err := r.db.Save(&dto).Error; err != nil {
		return nil, err
	}
	return ToMemberEntity(dto), nil
}

// Delete ...
func (r MemberRepository) Delete(member *entity.Member) error {
	dto := ToMemberDTO(member)
	if err := r.db.Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}
