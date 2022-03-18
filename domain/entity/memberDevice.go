package entity

import (
	"time"
)

// MemberDevice maps member and firebase device token ID
type MemberDevice struct {
	ID          uint
	MemberID    uint
	DeviceToken string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemberDevices ...
type MemberDevices []MemberDevice

// NewMemberDevice ...
func NewMemberDevice(memberID uint, token string) (*MemberDevice, error) {
	return &MemberDevice{
		MemberID:    memberID,
		DeviceToken: token,
	}, nil
}

// IsEqual guarantees Entity's identity
func (md *MemberDevice) IsEqual(other *MemberDevice) bool {
	return other.ID == md.ID
}
