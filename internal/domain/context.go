package domain

type ContextRequest struct {
	ThreadID  string
	ChannelID string
	Limit     int
}

type ContextStatus string

const (
	ContextStatusAllowed ContextStatus = "allowed"
	ContextStatusDenied  ContextStatus = "denied"
)

type ContextDenialCode string

const (
	ContextDenialNotAllowed  ContextDenialCode = "not_allowed"
	ContextDenialNotFound    ContextDenialCode = "not_found"
	ContextDenialUnsupported ContextDenialCode = "unsupported"
)

type ContextResponse struct {
	Status   ContextStatus
	Messages []Message
	Code     ContextDenialCode
}
