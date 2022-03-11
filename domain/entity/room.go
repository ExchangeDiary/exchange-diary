package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain"
	"github.com/jinzhu/copier"
)

// Room ...
type Room struct {
	ID     uint
	Name   string
	Code   string
	Hint   string
	Theme  string
	Period uint8

	MasterID      uint   // TODO: Member
	TurnAccountID uint   // TODO: Member
	Orders        []uint // []Member.ID
	Members       *Members

	CreatedAt *time.Time
	UpdatedAt *time.Time
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
		return 0, fmt.Errorf("There is no room member")
	}
	r.Orders, accountID = domain.Remove(r.Orders, accountID)
	if accountID == 0 {
		return 0, fmt.Errorf("There is no matched accountID from room.Orders")
	}
	return accountID, nil
}

// ChangeMaster ...
func (r *Room) ChangeMaster() error {
	nextCandidateIdx := 1
	if _, err := r.RemoveMember(r.MasterID); err != nil {
		return err
	}
	// Members에서 [] 가장 위에 존재하는 account id로 선출
	r.MasterID = (*r.Members)[nextCandidateIdx].ID
	return nil
}

// OrdersToJSON 는 []uint 타입을 []byte json타입으로 마샬링한다.
func (r *Room) OrdersToJSON() ([]byte, error) {
	orderJSON, err := json.Marshal(r.Orders)
	if err != nil {
		return nil, err
	}
	return []byte(orderJSON), nil
}

// MemberOnlyOrders returns master excluded memberIDs
func (r *Room) MemberOnlyOrders() ([]uint, error) {
	var orders []uint
	copier.Copy(&orders, &r.Orders)

	for i, accountID := range orders {
		if r.IsMaster(accountID) {
			return append(orders[:i], orders[i+1:]...), nil
		}
	}
	return nil, fmt.Errorf("There is no masterID in orders: %+v", r)
}

// NextTurn set room.TurnAccountID to next-turnAccountID and return it.
func (r *Room) NextTurn() (nextTurnAccountID uint) {
	curTurnAccountID := r.TurnAccountID

	if len(r.Orders) == 1 {
		nextTurnAccountID = curTurnAccountID
		return
	}

	var orders []uint
	copier.Copy(&orders, &r.Orders)
	for i, accountID := range orders {
		if accountID == curTurnAccountID {
			if i == len(r.Orders)-1 {
				nextTurnAccountID = orders[0]
			} else {
				nextTurnAccountID = orders[i+1]
			}
			break
		}
	}
	// set entity new turnAccountID
	r.TurnAccountID = nextTurnAccountID
	return
}

// NextTurnAt returns room.CreatedAt + period timestamp
func (r *Room) NextTurnAt() *time.Time {
	nt := r.CreatedAt.Add(time.Hour * 24 * time.Duration(r.Period))
	return &nt
}
