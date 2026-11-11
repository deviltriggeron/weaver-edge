package domain

type DecisionType string

const (
	DecisionTypeNoEvent           DecisionType = "NO_EVENT"
	DecisionTypeFinal             DecisionType = "FINAL"
	DecisionTypeNeedContext       DecisionType = "NEED_CONTEXT"
	DecisionTypeNeedClarification DecisionType = "NEED_CLARIFICATION"
)

type Decision struct {
	Type     DecisionType `json:"type"`
	Event    *Event       `json:"event,omitempty"`
	Reason   string       `json:"reason,omitempty"`
	Question string       `json:"question,omitempty"`
}

type Event struct {
	Type      string `json:"type"`
	Actor     string `json:"actor,omitempty"`
	Assignee  string `json:"assignee,omitempty"`
	Action    string `json:"action"`
	Object    string `json:"object,omitempty"`
	Deadline  string `json:"deadline,omitempty"`
	Reference string `json:"reference,omitempty"`
	Certainty string `json:"certainty,omitempty"`
}

type ClarificationReason string

const (
	ClarificationReasonAmbiguousReference ClarificationReason = "AMBIGUOUS_REFERENCE"
)
