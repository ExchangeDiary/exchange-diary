package vo

import "fmt"

// AlarmBody represents alarm body
type AlarmBody struct {
	RoomID     uint
	Code       string
	RoomName   string
	DiaryTitle string
	NickName   string
}

// NewAlarmBody ...
func NewAlarmBody(roomID uint, code TaskCode, roomName, diaryTitle, nickname string) *AlarmBody {
	return &AlarmBody{
		RoomID:     roomID,
		Code:       string(code),
		RoomName:   roomName,
		DiaryTitle: diaryTitle,
		NickName:   nickname,
	}
}

// ConvertToMap converts AlarmBody to map type
func (ab *AlarmBody) ConvertToMap() map[string]string {
	return map[string]string{
		"code":       ab.Code,
		"roomID":     fmt.Sprint(ab.RoomID),
		"roomName":   ab.RoomName,
		"diaryTitle": ab.DiaryTitle,
		"nickname":   ab.NickName,
	}
}
