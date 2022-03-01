package repository

import "github.com/ExchangeDiary/exchange-diary/domain/entity"

// MemberRepository ...
type MemberRepository interface {
	Create(member *entity.Member) (*entity.Member, error)
	GetByEmail(email string) (*entity.Member, error)
	Update(member *entity.Member) (*entity.Member, error)
	Delete(member *entity.Member) error
}
