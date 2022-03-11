package service

import (
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/tasks"

	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const (
	oneHour  = time.Hour * 1
	fourHour = time.Hour * 4
)

var rightNow = time.Time{}

// TaskService ...
type TaskService interface {
	DoRoomPeriodFINTask(roomID uint, baseURL string) error
	DoMemberOnDutyTask(email, deviceToken, baseURL string) error
	DoMemberBeforeTask(email, deviceToken, baseURL string, delta uint) error
	DoMemberPostedDiaryTask(roomID uint, deviceToken, baseURL string) error
}

type taskService struct {
	alarmService  AlarmService
	roomService   RoomService
	memberService MemberService
}

// NewTaskService ...
func NewTaskService(as AlarmService, rs RoomService, ms MemberService) TaskService {
	return &taskService{
		alarmService:  as,
		roomService:   rs,
		memberService: ms,
	}
}

func (ts *taskService) DoRoomPeriodFINTask(roomID uint, baseURL string) error {
	// 1. room update
	// 만약 방이 사라졌다면 error fin
	room, err := ts.roomService.Get(roomID)
	if err != nil {
		return err
	}

	turnAt := room.NextTurnAt()
	nxtTurnAccountID := room.NextTurn()
	if _, err := ts.roomService.Update(room); err != nil {
		return err
	}

	nxtMember, err := ts.memberService.Get(nxtTurnAccountID)
	if err != nil {
		return err
	}

	taskClient := tasks.GetClient()
	// 2. MEMBER_ON_DUTY task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberOnDutyCode).Encode(),
			taskspb.HttpMethod_POST,
			rightNow,
		)); err != nil {
		return err
	}

	// 3. MEMBER_BEFORE_1HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberBefore1HRCode).Encode(),
			taskspb.HttpMethod_POST,
			turnAt.Add(-oneHour))); err != nil {
		return err
	}
	// 4. MEMBER_BEFORE_4HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberBefore4HRCode).Encode(),
			taskspb.HttpMethod_POST,
			turnAt.Add(-fourHour))); err != nil {
		return err
	}
	// 5. Next ROOM_PERIOD_FIN task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			entity.NewTaskVO(roomID, nxtMember.Email, entity.RoomPeriodFinCode).Encode(),
			taskspb.HttpMethod_POST,
			*turnAt)); err != nil {
		return err
	}
	return nil
}

func (ts *taskService) DoMemberOnDutyTask(email, deviceToken, baseURL string) error {
	// 1. alarm to email, deviceToken user
	return nil
}

func (ts *taskService) DoMemberBeforeTask(email, deviceToken, baseURL string, delta uint) error {
	// 1. alarm to email, deviceToken user (by taskCode type)
	return nil
}

func (ts *taskService) DoMemberPostedDiaryTask(roomID uint, deviceToken, baseURL string) error {
	// 1. BroadCast alarm to RoomMember (except current member)
	// 2. 기존에 존재하는 ROOM_PERIOD_FIN task 업데이트 (바로 실행되도록 트리거)
	return nil
}
