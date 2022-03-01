package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain"
)

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
	Orders        []uint // []Member.ID

	CreatedAt time.Time
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

// IsAlreadyJoined determines whether account is master or member of room
func (r *Room) IsAlreadyJoined(accountID uint) bool {
	return r.IsMaster(accountID) || domain.Contains(r.Orders, accountID)
}

// AppendMember ...
func (r *Room) AppendMember(accountID uint) {
	r.Orders = append(r.Orders, accountID)
}

// RemoveMember ...
func (r *Room) RemoveMember(accountID uint) (uint, error) {
	if len(r.Orders) == 0 {
		return 0, errors.New("There is no room member")
	}
	r.Orders, accountID = domain.Remove(r.Orders, accountID)
	if accountID == 0 {
		return 0, errors.New("There is no matched accountID from room.Orders")
	}
	return accountID, nil
}

// ChangeMaster ...
func (r *Room) ChangeMaster() error {
	if _, err := r.RemoveMember(r.MasterID); err != nil {
		return err
	}
	// Order에서 가장 위에 존재하는 account id로 선출
	r.MasterID = r.Orders[0]
	return nil
}

// OrdersToJSON 는 []uint 타입을 []byte json타입으로 마샬링 변환한다.
func (r *Room) OrdersToJSON() ([]byte, error) {
	orderJSON, err := json.Marshal(r.Orders)
	if err != nil {
		return nil, err
	}
	return []byte(orderJSON), nil
}
