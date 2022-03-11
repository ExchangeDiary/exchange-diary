package service

import "github.com/ExchangeDiary/exchange-diary/domain/repository"

// TaskCode ...
type TaskCode string

const (
	// RoomPeriodFin task code type
	RoomPeriodFin TaskCode = "ROOM_PERIOD_FIN"
	// MemberOnDuty task code type
	MemberOnDuty = "MEMBER_ON_DUTY"
	// MemberBefore1HR task code type
	MemberBefore1HR = "MEMBER_BEFORE_1HR"
	// MemberBefore4HR task code type
	MemberBefore4HR = "MEMBER_BEFORE_4HR"
	// MemberPostedDiary task code type
	MemberPostedDiary = "MEMBER_POSTED_DIARY"
)

// TaskService ...
type TaskService interface {
	DoRoomPeriodFINTask(roomID uint, email, deviceToken string) error
	DoMemberOnDutyTask(email, deviceToken string) error
	DoMemberBeforeTask(email, deviceToken string, delta uint) error
	DoMemberPostedDiaryTask(roomID uint, deviceToken string) error
}

type taskService struct {
	alarmService   AlarmService
	roomRepository repository.RoomRepository
}

// NewTaskService ...
func NewTaskService(as AlarmService, rr repository.RoomRepository) TaskService {
	return &taskService{
		alarmService:   as,
		roomRepository: rr,
	}
}

func (ts *taskService) DoRoomPeriodFINTask(roomID uint, email, deviceToken string) error {
	// 1. room update // 만약 방이 사라졌다면 error fin
	// update 요소: nextTargetId
	// 2. MEMBER_ON_DUTY task register
	// 3. MEMBER_BEFORE_1HR task register
	// 3. MEMBER_BEFORE_4HR task register
	// 4. Next ROOM_PERIOD_FIN task register
	return nil
}

func (ts *taskService) DoMemberOnDutyTask(email, deviceToken string) error {
	// 1. alarm to email, deviceToken user
	return nil
}

func (ts *taskService) DoMemberBeforeTask(email, deviceToken string, delta uint) error {
	// 1. alarm to email, deviceToken user (by taskCode type)
	return nil
}

func (ts *taskService) DoMemberPostedDiaryTask(roomID uint, deviceToken string) error {
	// 1. BroadCast alarm to RoomMember (except current member)
	// 2. 기존에 존재하는 ROOM_PERIOD_FIN task 업데이트 (바로 실행되도록 트리거)
	return nil
}
