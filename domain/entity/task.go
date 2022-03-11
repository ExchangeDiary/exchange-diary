package entity

import (
	"bytes"
	"encoding/json"
)

// TaskVO ...
type TaskVO struct {
	RoomID uint
	Email  string
	Code   TaskCode
}

// TaskCode ...
type TaskCode string

const (
	// RoomPeriodFinCode task code type
	RoomPeriodFinCode TaskCode = "ROOM_PERIOD_FIN"
	// MemberOnDutyCode task code type
	MemberOnDutyCode = "MEMBER_ON_DUTY"
	// MemberBefore1HRCode task code type
	MemberBefore1HRCode = "MEMBER_BEFORE_1HR"
	// MemberBefore4HRCode task code type
	MemberBefore4HRCode = "MEMBER_BEFORE_4HR"
	// MemberPostedDiaryCode task code type
	MemberPostedDiaryCode = "MEMBER_POSTED_DIARY"
)

// NewTaskVO ...
func NewTaskVO(roomID uint, email string, code TaskCode) TaskVO {
	return TaskVO{
		RoomID: roomID,
		Email:  email,
		Code:   code,
	}
}

// Encode converts ValueObject to []byte
func (tv TaskVO) Encode() []byte {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(tv)
	return reqBodyBytes.Bytes()
}
