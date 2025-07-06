package domain

import (
	"github.com/google/uuid"
	"time"
)

type GiveawayPlanning struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	GiveawayID uuid.UUID              `json:"giveaway_id" db:"giveaway_id"`
	Date       time.Time              `json:"date" db:"date"`
	Type       GiveawayPlanningType   `json:"type" db:"type"`
	Status     GiveawayPlanningStatus `json:"status" db:"status"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" db:"updated_at"`
}

type GiveawayPlanningType string

const (
	GiveawayPlanningTypePublish GiveawayPlanningType = "PUBLISH"
	GiveawayPlanningTypeResult  GiveawayPlanningType = "RESULT"
)

type GiveawayPlanningStatus string

const (
	GiveawayPlanningStatusPending    GiveawayPlanningStatus = "PENDING"
	GiveawayPlanningStatusInProgress GiveawayPlanningStatus = "IN_PROGRESS"
	GiveawayPlanningStatusSuccess    GiveawayPlanningStatus = "SUCCESS"
	GiveawayPlanningStatusFailed     GiveawayPlanningStatus = "FAILED"
)
