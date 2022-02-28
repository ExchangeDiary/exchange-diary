package entity

import (
	"time"
)

type Member struct {
	ID                int
	Email             string
	Name              string
	ProfileURL        string
	AuthType          string
	TurnAlarmFlag     bool
	ActivityAlarmFlag bool
	CreatedAt         time.Time
}

type Members []Member

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
