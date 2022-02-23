package entity

type Room struct {
	ID    int
	Name  string
	Code  string
	Hint  string
	Theme string
}

type Rooms []Room

func NewRoom(name, code, hint, theme string) (*Room, error) {
	// TODO: field validation

	return &Room{
		Name:  name,
		Code:  code,
		Hint:  hint,
		Theme: theme,
	}, nil
}
