package dto

import "weaver-edge/internal/domain"

func MapSendMessageRequest(msg domain.Message, contextMessages []domain.Message) SendMessageRequest {
	contextDTO := make([]MessageRequest, 0, len(contextMessages))

	for _, contextMessage := range contextMessages {
		contextDTO = append(
			contextDTO,
			MapMessageToDTO(contextMessage),
		)
	}

	return SendMessageRequest{
		Message: MapMessageToDTO(msg),
		Context: contextDTO,
	}
}

func MapMessageToDTO(msg domain.Message) MessageRequest {
	return MessageRequest{
		ID:         msg.ID.String(),
		ChannelID:  msg.ChannelID,
		ThreadID:   msg.ThreadID,
		SenderID:   msg.SenderID,
		Text:       msg.Text,
		CreatedAt:  msg.CreatedAt,
		ReceivedAt: msg.ReceivedAt,
	}
}

func MapEventResponseToDecision(event EventResponse) domain.Decision {
	decision := domain.Decision{
		Type:     domain.DecisionType(event.Type),
		Reason:   string(event.Reason),
		Question: event.Question,
	}

	if event.Event != nil {
		decision.Event = &domain.Event{
			Type:      event.Event.Type,
			Actor:     event.Event.Actor,
			Assignee:  event.Event.Assignee,
			Action:    event.Event.Action,
			Object:    event.Event.Object,
			Deadline:  event.Event.Deadline,
			Reference: event.Event.Reference,
			Certainty: event.Event.Certainty,
		}
	}

	return decision
}
