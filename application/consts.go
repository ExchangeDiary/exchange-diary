package application

const (
	// CurrentMemberKey is current member mapping key which is set on gin.Context
	CurrentMemberKey = "currentMember"
	// AuthCodeKey is a auth code url query parameter
	AuthCodeKey = "authcode"
)

const (
	EntityNotFoundErr              = 400001
	EmptyParameterErr              = 400002
	InvalidRequestBodyErr          = 400003
	CannotCreateErr                = 400004
	TaskRoomPeriodFINCreateErr     = 400005
	TaskMemberOnDutyCreateErr      = 400006
	TaskMemberBefore1HRCreateErr   = 400007
	TaskMemberBefore4HRCreateErr   = 400008
	TaskMemberPostedDiaryCreateErr = 400009

	OnlyMemberErr         = 401001
	OnlyMasterErr         = 401002
	OnlyMemberOrMasterErr = 401003

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
