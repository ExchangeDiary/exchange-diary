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

	Master      Account
	TurnAccount Account
	Orders      []uint // []Account.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Rooms ...
type Rooms []Room

// NewRoom ...
func NewRoom(name, code, hint, theme string) (*Room, error) {
	// TODO: field validation

	return &Room{
		Name:  name,
		Code:  code,
		Hint:  hint,
		Theme: theme,
	}, nil
}
