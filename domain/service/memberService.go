package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
)

// MemberService ...
type MemberService interface {
	Create(email string, name string, profileURL string, authType string) (*entity.Member, error)
	Get(ID uint) (*entity.Member, error)
	GetByEmail(email string) (*entity.Member, error)
	Update(member *entity.Member) (*entity.Member, error)
	Delete(email string) error
	VerifyNickName(name string) (bool, error)
}

type memberService struct {
	memberRepository repository.MemberRepository
	memberService    MemberService
}

// NewMemberService ...
func NewMemberService(repository repository.MemberRepository) MemberService {
	return &memberService{
		memberRepository: repository,
	}
}

func (s *memberService) Create(email string, name string, profileURL string, authType string) (*entity.Member, error) {
	member, err := entity.NewMember(email, name, profileURL, authType)
	if err != nil {
		return nil, err
	}
	newMember, err := s.memberRepository.Create(member)
	if err != nil {
		return nil, err
	}
	return newMember, nil
}

func (s *memberService) Get(id uint) (*entity.Member, error) {
	member, err := s.memberRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) GetByEmail(email string) (*entity.Member, error) {
	member, err := s.memberRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) Update(member *entity.Member) (*entity.Member, error) {
	updatedMember, err := s.memberRepository.Update(member)
	if err != nil {
		return nil, err
	}
	return updatedMember, err
}

func (s memberService) Delete(email string) error {
	member, err := s.memberRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	err = s.memberRepository.Delete(member)
	if err != nil {
		return err
	}
	return nil
}

func (s memberService) VerifyNickName(name string) (bool, error) {

	member, err := s.memberRepository.GetByNickName(name)
	if err != nil {
		logger.Info(err.Error())
		return false, err
	}
	return member.IsNil(), nil
}
