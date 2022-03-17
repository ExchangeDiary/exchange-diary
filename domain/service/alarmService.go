package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/firebase"
)

// AlarmService ...
type AlarmService interface {
	PushByID(memberID uint, alarmBody *vo.AlarmBody) error
	PushByEmail(email string, alarmBody *vo.AlarmBody) error
	BroadCast(memberIDs []uint, alarmBody *vo.AlarmBody) error
}

type alarmService struct {
	memberService          MemberService
	memberDeviceRepository repository.MemberDeviceRepository
}

// NewAlarmService ...
func NewAlarmService(ms MemberService, mdr repository.MemberDeviceRepository) AlarmService {
	return &alarmService{
		memberService:          ms,
		memberDeviceRepository: mdr,
	}
}
func (as *alarmService) PushByID(memberID uint, alarmBody *vo.AlarmBody) (err error) {
	var deviceTokens []string
	if deviceTokens, err = as.memberDeviceRepository.GetAllTokens(memberID); err != nil {
		return
	}

	var failedTokens []string
	if failedTokens, err = firebase.GetClient().Push(deviceTokens, alarmBody); err != nil {
		return
	}

	if err = as.memberDeviceRepository.DeleteBatch(failedTokens); err != nil {
		return
	}
	return
}

func (as *alarmService) PushByEmail(email string, alarmBody *vo.AlarmBody) (err error) {
	var member *entity.Member
	if member, err = as.memberService.GetByEmail(email); err != nil {
		return
	}
	return as.PushByID(member.ID, alarmBody)
}

func (as *alarmService) BroadCast(memberIDs []uint, alarmBody *vo.AlarmBody) error {
	return nil
}
