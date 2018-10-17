package dsl

import "encoding/json"

// APIResponse is a response from the Telegram API with the result
// stored raw.
type dslTelegramBotAPIResponse struct {
	Ok          bool                              `json:"ok"`
	Result      json.RawMessage                   `json:"result"`
	ErrorCode   int                               `json:"error_code"`
	Description string                            `json:"description"`
	Parameters  *dslTelegramBotResponseParameters `json:"parameters"`
}

// ResponseParameters are various errors that can be returned in APIResponse.
type dslTelegramBotResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
	RetryAfter      int   `json:"retry_after"`        // optional
}

// User is a user on Telegram.
type dslTelegramBotUser struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
	IsBot        bool   `json:"is_bot"`        // optional
}

// Chat contains information about the place a message was sent.
type dslTelegramBotChat struct {
	ID                  int64  `json:"id"`
	Type                string `json:"type"`
	Title               string `json:"title"`                          // optional
	UserName            string `json:"username"`                       // optional
	FirstName           string `json:"first_name"`                     // optional
	LastName            string `json:"last_name"`                      // optional
	AllMembersAreAdmins bool   `json:"all_members_are_administrators"` // optional
	Description         string `json:"description,omitempty"`          // optional
	InviteLink          string `json:"invite_link,omitempty"`          // optional
}

// MessageEntity contains information about data in a Message.
type dslTelegramBotMessageEntity struct {
	Type   string              `json:"type"`
	Offset int                 `json:"offset"`
	Length int                 `json:"length"`
	URL    string              `json:"url"`  // optional
	User   *dslTelegramBotUser `json:"user"` // optional
}

// Contact contains information about a contact.
//
// Note that LastName and UserID may be empty.
type dslTelegramBotContact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"` // optional
	UserID      int    `json:"user_id"`   // optional
}

// Message is returned by almost every request, and contains data about
// almost anything.
type dslTelegramBotMessage struct {
	MessageID             int                            `json:"message_id"`
	From                  *dslTelegramBotUser            `json:"from"` // optional
	Date                  int                            `json:"date"`
	Chat                  *dslTelegramBotChat            `json:"chat"`
	ForwardFrom           *dslTelegramBotUser            `json:"forward_from"`            // optional
	ForwardFromChat       *dslTelegramBotChat            `json:"forward_from_chat"`       // optional
	ForwardFromMessageID  int                            `json:"forward_from_message_id"` // optional
	ForwardDate           int                            `json:"forward_date"`            // optional
	ReplyToMessage        *dslTelegramBotMessage         `json:"reply_to_message"`        // optional
	EditDate              int                            `json:"edit_date"`               // optional
	Text                  string                         `json:"text"`                    // optional
	Entities              *[]dslTelegramBotMessageEntity `json:"entities"`                // optional
	Caption               string                         `json:"caption"`                 // optional
	Contact               *dslTelegramBotContact         `json:"contact"`                 // optional
	NewChatMembers        *[]dslTelegramBotUser          `json:"new_chat_members"`        // optional
	LeftChatMember        *dslTelegramBotUser            `json:"left_chat_member"`        // optional
	NewChatTitle          string                         `json:"new_chat_title"`          // optional
	DeleteChatPhoto       bool                           `json:"delete_chat_photo"`       // optional
	GroupChatCreated      bool                           `json:"group_chat_created"`      // optional
	SuperGroupChatCreated bool                           `json:"supergroup_chat_created"` // optional
	ChannelChatCreated    bool                           `json:"channel_chat_created"`    // optional
	MigrateToChatID       int64                          `json:"migrate_to_chat_id"`      // optional
	MigrateFromChatID     int64                          `json:"migrate_from_chat_id"`    // optional
	PinnedMessage         *dslTelegramBotMessage         `json:"pinned_message"`          // optional
}

// Update is an update response, from GetUpdates.
type dslTelegramBotUpdate struct {
	UpdateID          int                          `json:"update_id"`
	Message           *dslTelegramBotMessage       `json:"message"`
	EditedMessage     *dslTelegramBotMessage       `json:"edited_message"`
	ChannelPost       *dslTelegramBotMessage       `json:"channel_post"`
	EditedChannelPost *dslTelegramBotMessage       `json:"edited_channel_post"`
	CallbackQuery     *dslTelegramBotCallbackQuery `json:"callback_query"`
}

type dslTelegramBotCallbackQuery struct {
	ID              string                 `json:"id"`
	From            *dslTelegramBotUser    `json:"from"`
	Message         *dslTelegramBotMessage `json:"message"`           // optional
	InlineMessageID string                 `json:"inline_message_id"` // optional
	ChatInstance    string                 `json:"chat_instance"`
	Data            string                 `json:"data"` // optional
}
