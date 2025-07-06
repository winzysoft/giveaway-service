package domain

import (
	"github.com/google/uuid"
	"time"
)

type Giveaway struct {
	Id                uuid.UUID      `json:"id" db:"id"`
	Status            GiveawayStatus `json:"status" db:"status"`
	ChatId            string         `json:"chatId" db:"chat_id"`
	MessageId         string         `json:"messageId" db:"message_id"`
	ShouldEditMessage bool           `json:"shouldEditMessage" db:"should_edit_message"`
	ShouldSendNew     bool           `json:"shouldSendNew" db:"should_send_new"`
	WinCount          int            `json:"winCount" db:"win_count"`
	CreatedAt         time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time      `json:"updatedAt" db:"updated_at"`
}

type GiveawayStatus string

// Возможные значения GiveawayStatus
const (
	GiveawayStatusCreated  GiveawayStatus = "CREATED"
	GiveawayStatusActive   GiveawayStatus = "ACTIVE"
	GiveawayStatusDraft    GiveawayStatus = "DRAFT"
	GiveawayStatusEnded    GiveawayStatus = "ENDED"
	GiveawayStatusCanceled GiveawayStatus = "CANCELED"
	GiveawayStatusDeleted  GiveawayStatus = "DELETED"
)
