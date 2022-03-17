package service

import (
	"fmt"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/vo"
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
	DoMemberOnDutyTask(email string) (err error)
	DoMemberBeforeTask(email string, code vo.TaskCode) (err error)
	DoMemberPostedDiaryTask(roomID uint, baseURL string) error

	RegisterRoomPeriodFINTask(c *tasks.Client, baseURL string, roomID, accountID uint, dueAt *time.Time) (taskID string, err error)
	RegisterMemberPostedDiaryTask(roomID uint, baseURL string) (taskID string, err error)

	GetTask(c *tasks.Client, code vo.TaskCode, roomID, turnAccountID uint) (*taskspb.Task, error)
	DeleteTask(c *tasks.Client, code vo.TaskCode, roomID, turnAccountID uint) error
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
			genUniqueTaskID(roomID, nxtTurnAccountID, vo.MemberOnDutyCode),
			vo.NewTaskVO(roomID, nxtMember.Email, vo.MemberOnDutyCode).Encode(),
			taskspb.HttpMethod_POST,
			rightNow,
		)); err != nil {
		return err
	}

	// 3. MEMBER_BEFORE_1HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			genUniqueTaskID(roomID, nxtTurnAccountID, vo.MemberBefore1HRCode),
			vo.NewTaskVO(roomID, nxtMember.Email, vo.MemberBefore1HRCode).Encode(),
			taskspb.HttpMethod_POST,
			dueAt.Add(-oneHour))); err != nil {
		return err
	}
	// 4. MEMBER_BEFORE_4HR task register
	if _, err := taskClient.RegisterTask(
		taskClient.BuildTask(baseURL,
			genUniqueTaskID(roomID, nxtTurnAccountID, vo.MemberBefore4HRCode),
			vo.NewTaskVO(roomID, nxtMember.Email, vo.MemberBefore4HRCode).Encode(),
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
	code := vo.RoomPeriodFinCode
	if task, err = c.RegisterTask(
		c.BuildTask(
			baseURL,
			genUniqueTaskID(roomID, accountID, code),
			vo.NewTaskVO(roomID, "", code).Encode(),
			taskspb.HttpMethod_POST,
			*dueAt,
		),
	); err != nil {
		return "", err
	}
	return task.Name, nil
}

func (ts *taskService) DoMemberOnDutyTask(email string) (err error) {
	if err = ts.alarmService.PushByEmail(email, vo.MemberOnDutyCode); err != nil {
		return
	}
	return
}

func (ts *taskService) DoMemberBeforeTask(email string, code vo.TaskCode) (err error) {
	if err = ts.alarmService.PushByEmail(email, code); err != nil {
		return
	}
	return
}

func (ts *taskService) DoMemberPostedDiaryTask(roomID uint, baseURL string) error {
	room, err := ts.roomService.Get(roomID)
	if err != nil {
		return err
	}

	// 1. BroadCast alarm to RoomMember (except current member)
	if err = ts.alarmService.BroadCast(room.MemberAllExceptTurnAccount(), vo.MemberPostedDiaryCode); err != nil {
		return err
	}

	// 2. 기존에 존재하는 ROOM_PERIOD_FIN task 업데이트 (바로 실행되도록 트리거)
	taskClient := tasks.GetClient()
	taskID := genUniqueTaskID(room.ID, room.TurnAccountID, vo.RoomPeriodFinCode)
	if err := taskClient.DeleteTask(taskID); err != nil {
		return err
	}
	return ts.DoRoomPeriodFINTask(room.ID, baseURL)
}

func (ts *taskService) RegisterMemberPostedDiaryTask(roomID uint, baseURL string) (taskID string, err error) {
	// TODO: diary crud 만들어지면 구현
	return "", nil
}

func (ts *taskService) UpdateRoomPeriodFINTaskETA(c *tasks.Client, baseURL string, roomID, turnAccountID uint, nextDueAt *time.Time) error {
	code := vo.RoomPeriodFinCode
	taskID := genUniqueTaskID(roomID, turnAccountID, vo.RoomPeriodFinCode)
	if _, err := c.UpdateTask(
		taskID,
		c.BuildTask(baseURL,
			genUniqueTaskID(roomID, turnAccountID, code),
			vo.NewTaskVO(roomID, "", code).Encode(),
			taskspb.HttpMethod_POST,
			*nextDueAt,
		)); err != nil {
		return err
	}
	return nil
}

func (ts *taskService) GetTask(c *tasks.Client, code vo.TaskCode, roomID, turnAccountID uint) (task *taskspb.Task, err error) {
	taskID := genUniqueTaskID(roomID, turnAccountID, code)
	if task, err = c.GetTask(taskID); err != nil {
		return nil, err
	}
	return
}

func (ts *taskService) DeleteTask(c *tasks.Client, code vo.TaskCode, roomID, turnAccountID uint) error {
	taskID := genUniqueTaskID(roomID, turnAccountID, code)
	if err := c.DeleteTask(taskID); err != nil {
		return err
	}
	return nil
}

// room_id + AccountID + eventCode
func genUniqueTaskID(roomID, turnAccountID uint, code vo.TaskCode) string {
	return fmt.Sprintf("%d-%d-%s", roomID, turnAccountID, code)
}
