package service

import (
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/firebase"
)

// AlarmService ...
type AlarmService interface {
	Create(memberID, roomID uint, code vo.TaskCode, roomName, diaryTitle, authorNickname string) (*entity.Alarm, error)
	GetAll(memberID uint) (*entity.Alarms, error)
	PushByID(memberID uint, al *entity.Alarm) error
	PushByEmail(email string, al *entity.Alarm) error
	BroadCast(memberIDs []uint, al *entity.Alarm) error
}

type alarmService struct {
	memberService          MemberService
	memberDeviceRepository repository.MemberDeviceRepository
	alarmRepository        repository.AlarmRepository
}

// NewAlarmService ...
func NewAlarmService(ms MemberService, mdr repository.MemberDeviceRepository, ar repository.AlarmRepository) AlarmService {
	return &alarmService{
		memberService:          ms,
		memberDeviceRepository: mdr,
		alarmRepository:        ar,
	}
}

func (as *alarmService) Create(memberID, roomID uint, code vo.TaskCode, roomName, diaryTitle, authorNickname string) (*entity.Alarm, error) {
	alarm, err := as.alarmRepository.Create(entity.NewAlarm(memberID, roomID, code, roomName, diaryTitle, authorNickname))
	if err != nil {
		return nil, err
	}
	return alarm, nil
}

func (as *alarmService) GetAll(memberID uint) (*entity.Alarms, error) {
	return as.alarmRepository.GetAll(memberID)
}

func (as *alarmService) PushByID(memberID uint, al *entity.Alarm) (err error) {
	var deviceTokens []string
	if deviceTokens, err = as.memberDeviceRepository.GetAllTokens(memberID); err != nil {
		return
	}

	var failedTokens []string
	if failedTokens, err = firebase.GetClient().Push(deviceTokens, al); err != nil {
		return
	}
	if len(failedTokens) > 0 {
		if err = as.memberDeviceRepository.DeleteBatch(failedTokens); err != nil {
			return
		}
	}

	return
}

func (as *alarmService) PushByEmail(email string, al *entity.Alarm) (err error) {
	var member *entity.Member
	if member, err = as.memberService.GetByEmail(email); err != nil {
		return
	}
	return as.PushByID(member.ID, al)
}

func (as *alarmService) BroadCast(memberIDs []uint, al *entity.Alarm) (err error) {
	var deviceTokens []string
	if deviceTokens, err = as.memberDeviceRepository.GetAllMemberTokens(memberIDs); err != nil {
		return
	}

	var failedTokens []string
	if failedTokens, err = firebase.GetClient().Push(deviceTokens, al); err != nil {
		return
	}
	if len(failedTokens) > 0 {
		if err = as.memberDeviceRepository.DeleteBatch(failedTokens); err != nil {
			return
		}
	}
	return
}
