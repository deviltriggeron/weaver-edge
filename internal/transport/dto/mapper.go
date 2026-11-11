package dto

import "weaver-edge/internal/domain"

func ContextReqToDomain(req ContextRequest) domain.ContextRequest {
	return domain.ContextRequest{
		ThreadID:  req.ThreadID,
		ChannelID: req.ChannelID,
		Limit:     req.Limit,
	}
}

func ContextRespDomainToResp(resp *domain.ContextResponse) ContextResponse {
	return ContextResponse{
		Status:   ContextStatus(resp.Status),
		Messages: messagesToDTO(resp.Messages),
		Code:     ContextDenialCode(resp.Code),
	}
}

func messagesToDTO(msg []domain.Message) []ContextMessage {
	var messages []ContextMessage

	for _, m := range msg {
		messages = append(messages, ContextMessage{
			ID:        m.ExternalID,
			ThreadID:  m.ThreadID,
			Text:      m.Text,
			SenderID:  m.SenderID,
			ChannelID: m.ChannelID,
			CreatedAt: m.CreatedAt,
		})
	}

	return messages
}
