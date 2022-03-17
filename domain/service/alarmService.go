package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/firebase"
)

// AlarmService ...
type AlarmService interface {
	PushByID(memberID uint, code vo.TaskCode) error
	PushByEmail(email string, code vo.TaskCode) error
	BroadCast(memberIDs []uint, code vo.TaskCode) error
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
func (as *alarmService) PushByID(memberID uint, code vo.TaskCode) (err error) {
	var deviceTokens []string
	if deviceTokens, err = as.memberDeviceRepository.GetAllTokens(memberID); err != nil {
		return
	}

	var failedTokens []string
	if failedTokens, err = firebase.GetClient().Push(deviceTokens, vo.NewAlarmBody(1, code, "", "", "")); err != nil {
		return
	}

	if err = as.memberDeviceRepository.DeleteBatch(failedTokens); err != nil {
		return
	}
	return
}

func (as *alarmService) PushByEmail(email string, code vo.TaskCode) (err error) {
	var member *entity.Member
	if member, err = as.memberService.GetByEmail(email); err != nil {
		return
	}
	return as.PushByID(member.ID, code)
}

func (as *alarmService) BroadCast(memberIDs []uint, code vo.TaskCode) error {
	return nil
}
