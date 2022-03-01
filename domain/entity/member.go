package entity

import (
	"time"
)

// Member ...
type Member struct {
	ID                uint
	Email             string
	Name              string
	ProfileURL        string
	AuthType          string
	TurnAlarmFlag     bool
	ActivityAlarmFlag bool
	CreatedAt         time.Time
}

// Members ...
type Members []Member

// NewMember ...
func NewMember(email, name, profileURL string, authType string) (*Member, error) {
	return &Member{
		Email:             email,
		Name:              name,
		ProfileURL:        profileURL,
		AuthType:          authType,
		TurnAlarmFlag:     true,
		ActivityAlarmFlag: true,
	}, nil
}

// IsEqual guarantees Entity's identity
func (a *Member) IsEqual(other *Member) bool {
	return (other.ID == a.ID) && (other.Email == a.Email)
}
