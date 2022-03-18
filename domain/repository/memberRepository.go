package repository

import "github.com/ExchangeDiary/exchange-diary/domain/entity"

// MemberRepository ...
type MemberRepository interface {
	Create(member *entity.Member) (*entity.Member, error)
	Get(id uint) (*entity.Member, error)
	GetByEmail(email string) (*entity.Member, error)
	GetAllByIDs(ids []uint) (*entity.Members, error)
	Update(member *entity.Member) (*entity.Member, error)
	Delete(member *entity.Member) error
}
