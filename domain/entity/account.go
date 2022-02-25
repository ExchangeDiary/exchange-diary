package entity

import "time"

// Account ...
type Account struct {
	ID                int
	Email             string
	Name              string
	ProfileURL        string
	TurnAlarmFlag     bool
	ActivityAlarmFlag bool
	CreatedAt         time.Time
}

// Accounts ...
type Accounts []Account

// NewAccount ...
func NewAccount(email, name, profileURL string) (*Account, error) {
	// TODO: field validation

	return &Account{
		Email:             email,
		Name:              name,
		ProfileURL:        profileURL,
		TurnAlarmFlag:     true,
		ActivityAlarmFlag: true,
	}, nil
}

// IsEqual guarantees Entity's identity
func (a *Account) IsEqual(other *Account) bool {
	return (other.ID == a.ID) && (other.Email == a.Email)
}
