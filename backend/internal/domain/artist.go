package domain

import (
	"time"

	"github.com/google/uuid"
)

// Artist is a scrape target whose official site is monitored for concerts.
type Artist struct {
	ID              uuid.UUID
	Name            string
	OfficialSiteURL string
	CreatedAt       time.Time
}
