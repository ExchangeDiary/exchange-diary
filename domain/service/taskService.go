package service

import (
	"fmt"
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

	RegisterRoomPeriodFINTask(c *tasks.Client, baseURL string, roomID, accountID uint, dueAt *time.Time) (taskID string, err error)
	RegisterMemberPostedDiaryTask(roomID uint, deviceToken, baseURL string) (taskID string, err error)

	UpdateRoomPeriodFINTaskETA(c *tasks.Client, baseURL string, roomID, turnAccountID uint, nextDueAt *time.Time) error
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

	dueAt := room.NextDueAt()
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
			genUniqueTaskID(roomID, nxtTurnAccountID, entity.MemberOnDutyCode),
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberOnDutyCode).Encode(),
			taskspb.HttpMethod_POST,
			rightNow,
		)); err != nil {
		return err
	}

	// 3. MEMBER_BEFORE_1HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			genUniqueTaskID(roomID, nxtTurnAccountID, entity.MemberBefore1HRCode),
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberBefore1HRCode).Encode(),
			taskspb.HttpMethod_POST,
			dueAt.Add(-oneHour))); err != nil {
		return err
	}
	// 4. MEMBER_BEFORE_4HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			genUniqueTaskID(roomID, nxtTurnAccountID, entity.MemberBefore4HRCode),
			entity.NewTaskVO(roomID, nxtMember.Email, entity.MemberBefore4HRCode).Encode(),
			taskspb.HttpMethod_POST,
			dueAt.Add(-fourHour))); err != nil {
		return err
	}
	// 5. Next ROOM_PERIOD_FIN task register
	if _, err := ts.RegisterRoomPeriodFINTask(
		taskClient,
		baseURL,
		roomID,
		nxtTurnAccountID,
		dueAt,
	); err != nil {
		return err
	}
	return nil
}

func (ts *taskService) RegisterRoomPeriodFINTask(c *tasks.Client, baseURL string, roomID, accountID uint, dueAt *time.Time) (taskID string, err error) {
	var task *taskspb.Task
	code := entity.RoomPeriodFinCode
	if task, err = c.RegisterTask(
		c.BuildTask(
			baseURL,
			genUniqueTaskID(roomID, accountID, code),
			entity.NewTaskVO(roomID, "", code).Encode(),
			taskspb.HttpMethod_POST,
			*dueAt,
		),
	); err != nil {
		return "", err
	}
	return task.Name, nil
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

func (ts *taskService) RegisterMemberPostedDiaryTask(roomID uint, deviceToken, baseURL string) (taskID string, err error) {
	// TODO:
	return "", nil
}

func (ts *taskService) UpdateRoomPeriodFINTaskETA(c *tasks.Client, baseURL string, roomID, turnAccountID uint, nextDueAt *time.Time) error {
	code := entity.RoomPeriodFinCode
	taskID := genUniqueTaskID(roomID, turnAccountID, entity.RoomPeriodFinCode)
	if _, err := c.UpdateTask(
		taskID,
		c.BuildTask(baseURL,
			genUniqueTaskID(roomID, turnAccountID, code),
			entity.NewTaskVO(roomID, "", code).Encode(),
			taskspb.HttpMethod_POST,
			*nextDueAt,
		)); err != nil {
		return err
	}
	return nil
}

// room_id + AccountID + eventCode
func genUniqueTaskID(roomID, turnAccountID uint, code entity.TaskCode) string {
	return fmt.Sprintf("%d-%d-%s", roomID, turnAccountID, code)
}
