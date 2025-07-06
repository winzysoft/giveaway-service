package dto

import "giveaway-service/internal/domain"

type GiveawaySaveRequest struct {
	Status            domain.GiveawayStatus `json:"status"` // todo: лучше разделить на dto и domain статусы
	ChatId            string                `json:"chatId"`
	ManagerIds        []string              `json:"managerIds"`
	MessageId         string                `json:"messageId"`
	ShouldEditMessage bool                  `json:"shouldEditMessage"`
	ShouldSendNew     bool                  `json:"shouldSendNew"`
	WinCount          int                   `json:"winCount"`
}
