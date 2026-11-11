package mapper

import (
	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/db/model"
)

func MsgModelToDomain(msg model.MessageModel) domain.Message {
	return domain.Message{
		ID:         msg.ID,
		ExternalID: msg.ExternalID,
		Platform:   domain.Platform(msg.Platform),
		ThreadID:   msg.ThreadID,
		ChannelID:  msg.ChannelID,
		SenderID:   msg.SenderID,
		Text:       msg.Text,
		CreatedAt:  msg.CreatedAt,
		ReceivedAt: msg.ReceivedAt,
	}
}

func MsgDomainToModel(msg domain.Message) model.MessageModel {
	return model.MessageModel{
		ID:         msg.ID,
		ExternalID: msg.ExternalID,
		Platform:   string(msg.Platform),
		ThreadID:   msg.ThreadID,
		ChannelID:  msg.ChannelID,
		SenderID:   msg.SenderID,
		Text:       msg.Text,
		CreatedAt:  msg.CreatedAt,
		ReceivedAt: msg.ReceivedAt,
	}
}
