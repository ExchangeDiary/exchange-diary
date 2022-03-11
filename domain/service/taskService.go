package service

import (
	"strconv"
	"time"

	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/tasks"

	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

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

const (
	oneHour  = time.Hour * 1
	fourHour = time.Hour * 4
)

// TaskService ...
type TaskService interface {
	DoRoomPeriodFINTask(roomID uint, email, deviceToken, baseURL string) error
	DoMemberOnDutyTask(email, deviceToken, baseURL string) error
	DoMemberBeforeTask(email, deviceToken, baseURL string, delta uint) error
	DoMemberPostedDiaryTask(roomID uint, deviceToken, baseURL string) error
}

type taskService struct {
	alarmService AlarmService
	roomService  RoomService
}

// NewTaskService ...
func NewTaskService(as AlarmService, rs RoomService) TaskService {
	return &taskService{
		alarmService: as,
		roomService:  rs,
	}
}

func (ts *taskService) DoRoomPeriodFINTask(roomID uint, email, deviceToken, baseURL string) error {
	// 1. room update
	// 만약 방이 사라졌다면 error fin
	room, err := ts.roomService.Get(roomID)
	if err != nil {
		return err
	}
	_ = room.NextTurn() // TODO: nxtTurnAccountID to OIDC && registerTask to member_id
	turnAt := room.NextTurnAt()
	if _, err := ts.roomService.Update(room); err != nil {
		return err
	}

	taskClient := tasks.GetClient()
	taskID := MemberOnDuty + "_ROOM_" + strconv.Itoa(int(roomID))
	// 2. MEMBER_ON_DUTY task register
	if _, err := taskClient.RunTask(taskClient.TaskID(taskID), baseURL, taskspb.HttpMethod_POST); err != nil {
		return err
	}
	// 3. MEMBER_BEFORE_1HR task register
	if _, err := taskClient.RegisterTask(baseURL, "", taskspb.HttpMethod_POST, turnAt.Add(-oneHour)); err != nil {
		return err
	}
	// 3. MEMBER_BEFORE_4HR task register
	if _, err := taskClient.RegisterTask(baseURL, "", taskspb.HttpMethod_POST, turnAt.Add(-fourHour)); err != nil {
		return err
	}
	// 4. Next ROOM_PERIOD_FIN task register
	if _, err := taskClient.RegisterTask(baseURL, "", taskspb.HttpMethod_POST, *turnAt); err != nil {
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
