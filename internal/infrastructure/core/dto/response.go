package dto

const (
	ReasonResultNoEvent           ReasonResultType = "NO_EVENT"
	ReasonResultFinal             ReasonResultType = "FINAL"
	ReasonResultNeedContext       ReasonResultType = "NEED_CONTEXT"
	ReasonResultNeedClarification ReasonResultType = "NEED_CLARIFICATION"
)

const (
	ClarificationAmbiguousReference ClarificationReason = "AMBIGUOUS_REFERENCE"
)

type (
	ReasonResultType string

	ClarificationReason string

	EventResponse struct { // ai reason -> core -> edge
		Type     ReasonResultType    `json:"type"`
		Event    *EventResponseEvent `json:"event,omitempty"`
		Reason   ClarificationReason `json:"reason,omitempty"`
		Question string              `json:"question,omitempty"`
		Request  *ContextRequest     `json:"request,omitempty"`
	}

	EventResponseEvent struct {
		Type      string `json:"type"`
		Actor     string `json:"actor"`
		Assignee  string `json:"assignee"`
		Action    string `json:"action"`
		Object    string `json:"object"`
		Deadline  string `json:"deadline"`
		Reference string `json:"reference"`
		Certainty string `json:"certainty"`
	}

	ClarificationResponse struct {
		ID       string `json:"id"`
		Reason   string `json:"reason"`
		Question string `json:"question"`
	}

	ContextRequest struct {
		ThreadID  string `json:"thread_id"`
		ChannelID string `json:"channel_id"`
		Scope     string `json:"scope"`
		Limit     int    `json:"limit"`
	}
)
