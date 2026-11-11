package dto

// delete
type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID       int64  `json:"message_id"`
	MessageThreadID int64  `json:"message_thread_id"`
	From            User   `json:"from"`
	Chat            Chat   `json:"chat"`
	Date            int64  `json:"date"`
	Text            string `json:"text,omitempty"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
}

// delete

type TelegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

type TelegramMessage struct {
	MessageID       int64        `json:"message_id"`
	MessageThreadID int64        `json:"message_thread_id"`
	From            TelegramUser `json:"from"`
	Chat            TelegramChat `json:"chat"`
	Date            int64        `json:"date"`
	Text            string       `json:"text,omitempty"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
}

type TelegramChat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	Type      string `json:"type"`
}
