package vo

import (
	"fmt"

	"github.com/ExchangeDiary/exchange-diary/domain"
)

// Alarm represents alarm body
type Alarm struct {
	RoomID         uint
	Code           string
	RoomName       string
	DiaryTitle     string
	AuthorNickName string
}

// NewAlarm ...
func NewAlarm(roomID uint, code TaskCode, roomName, diaryTitle, authorNickname string) *Alarm {
	return &Alarm{
		RoomID:         roomID,
		Code:           string(code),
		RoomName:       roomName,
		DiaryTitle:     diaryTitle,
		AuthorNickName: authorNickname,
	}
}

// ConvertToMap converts Alarm to map type
func (ab *Alarm) ConvertToMap() (alarmMap map[string]string) {
	now := domain.CurrentDateTime()
	switch ab.Code {
	case MemberOnDutyCode:
		alarmMap = map[string]string{
			"code":     ab.Code,
			"title":    "내가 일기 쓸 차례에요!",
			"roomName": ab.RoomName,
			"alarm_at": now.String(),
		}
	case MemberBefore1HRCode:
		alarmMap = map[string]string{
			"code":     ab.Code,
			"title":    "일기 등록까지 1시간 남았어요!",
			"roomName": ab.RoomName,
			"alarm_at": now.String(),
		}
	case MemberBefore4HRCode:
		alarmMap = map[string]string{
			"code":     ab.Code,
			"title":    "일기 등록까지 4시간 남았어요!",
			"roomName": ab.RoomName,
			"alarm_at": now.String(),
		}
	case MemberPostedDiaryCode:
		alarmMap = map[string]string{
			"code":     ab.Code,
			"title":    fmt.Sprintf("'%s' 새글 등록", ab.DiaryTitle),
			"roomName": ab.RoomName,
			"alarm_at": now.String(),
			"author":   ab.AuthorNickName,
		}
	}
	return alarmMap
}
