package dto

import (
	"github.com/google/uuid"
	"giveaway-service/internal/domain"
	"time"
)

type GiveawaySaveResponse struct {
	Id                uuid.UUID             `json:"id"` // todo: лучше не отдавать id из базы клиенту
	Status            domain.GiveawayStatus `json:"status"`
	MetaId            string                `json:"metaId"`
	ChatId            string                `json:"chatId"`
	MessageId         string                `json:"messageId"`
	ShouldEditMessage bool                  `json:"shouldEditMessage"`
	ShouldSendNew     bool                  `json:"shouldSendNew"`
	WinCount          int                   `json:"winCount"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}
