package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
)

// Alarm represents alarm body
type Alarm struct {
	MemberID uint
	RoomID   uint
	Code     string
	Title    string
	RoomName string
	Author   string
	AlarmAt  *time.Time
}

// Alarms ...
type Alarms []Alarm

// NewAlarm ...
func NewAlarm(memberID, roomID uint, code vo.TaskCode, roomName, diaryTitle, authorNickname string) *Alarm {
	now := domain.CurrentDateTime()
	switch code {
	case vo.MemberOnDutyCode:
		return &Alarm{
			MemberID: memberID,
			RoomID:   roomID,
			Code:     string(code),
			RoomName: roomName,
			AlarmAt:  &now,
			Title:    "내가 일기 쓸 차례에요!",
		}
	case vo.MemberBefore1HRCode:
		return &Alarm{
			MemberID: memberID,
			RoomID:   roomID,
			Code:     string(code),
			RoomName: roomName,
			AlarmAt:  &now,
			Title:    "일기 등록까지 1시간 남았어요!",
		}
	case vo.MemberBefore4HRCode:
		return &Alarm{
			MemberID: memberID,
			RoomID:   roomID,
			Code:     string(code),
			RoomName: roomName,
			AlarmAt:  &now,
			Title:    "일기 등록까지 4시간 남았어요!",
		}
	case vo.MemberPostedDiaryCode:
		return &Alarm{
			MemberID: memberID,
			RoomID:   roomID,
			Code:     string(code),
			RoomName: roomName,
			AlarmAt:  &now,
			Title:    fmt.Sprintf("'%s' 새글 등록", diaryTitle),
			Author:   authorNickname,
		}
	default:
		fmt.Printf("'%s' is invalid code type", code)
		return nil
	}
}

// UnqFields returns alarm's unique field component.
func (a *Alarm) UnqFields() (roomID, memberID uint, code vo.TaskCode) {
	return a.RoomID, a.MemberID, vo.TaskCode(a.Code)
}

// ToMap converts Alarm to map type
func (a *Alarm) ToMap() (alarmMap map[string]string) {
	raw, _ := json.Marshal(a)
	json.Unmarshal(raw, &alarmMap)
	return

}
