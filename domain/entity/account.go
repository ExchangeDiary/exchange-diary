package entity

// Account ...
type Account struct {
	ID   int
	Name string
}

// Accounts ...
type Accounts []Account

// NewAccount ...
func NewAccount(name string) (*Account, error) {
	// TODO: field validation

	return &Account{
		Name: name,
	}, nil
}
