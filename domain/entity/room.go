package entity

import "time"

// Room ...
type Room struct {
	ID     uint
	Name   string
	Code   string
	Hint   string
	Theme  string
	Period uint8

	MasterID      uint
	TurnAccountID uint
	Orders        []uint // []Account.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Rooms ...
type Rooms []Room

// NewRoom ...
func NewRoom(masterID uint, name, code, hint, theme string, period uint8) (*Room, error) {
	// TODO: field validation

	orders := []uint{masterID}
	return &Room{
		Name:          name,
		Code:          code,
		Hint:          hint,
		Theme:         theme,
		Period:        period,
		MasterID:      masterID,
		TurnAccountID: masterID,
		Orders:        orders,
	}, nil
}

// IsEqual guarantees Entity's identity
func (r *Room) IsEqual(other *Room) bool {
	return other.ID == r.ID
}

// IsMaster returns whether account is a room's master or not
func (r *Room) IsMaster(accountID uint) bool {
	return r.MasterID == accountID
}
