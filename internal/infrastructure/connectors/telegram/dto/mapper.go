package dto

import (
	"fmt"
	"strconv"
	"time"

	"weaver-edge/internal/domain"

	"github.com/google/uuid"
)

func MapMessage(src *TelegramMessage) domain.Message {
	now := time.Now().UTC()

	return domain.Message{
		ID:         uuid.New(),
		ExternalID: fmt.Sprintf("%d:%d", src.Chat.ID, src.MessageID),
		Platform:   domain.Telegram,
		ThreadID:   strconv.FormatInt(src.MessageThreadID, 10),
		Text:       src.Text,
		SenderID:   strconv.FormatInt(src.From.ID, 10),
		ChannelID:  strconv.FormatInt(src.Chat.ID, 10),
		CreatedAt:  time.Unix(src.Date, 0).UTC(),
		ReceivedAt: now,
	}
}
