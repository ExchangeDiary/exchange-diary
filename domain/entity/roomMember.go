package entity

import (
	"time"
)

// RoomMember maps room and accout
type RoomMember struct {
	ID        uint
	RoomID    uint
	AccountID uint
	CreatedAt time.Time
}

// RoomMembers ...
type RoomMembers []RoomMember

// NewRoomMember ...
func NewRoomMember(roomID, accountID uint) (*RoomMember, error) {
	return &RoomMember{
		RoomID:    roomID,
		AccountID: accountID,
	}, nil
}

// IsEqual guarantees Entity's identity
func (r *RoomMember) IsEqual(other *RoomMember) bool {
	return other.ID == r.ID
}
