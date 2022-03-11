package service

// AlarmService ...
type AlarmService interface {
}

type alarmService struct{}

// NewAlarmService ...
func NewAlarmService() AlarmService {
	return &alarmService{}
}
