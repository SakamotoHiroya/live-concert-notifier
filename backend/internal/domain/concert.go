package domain

import (
	"time"

	"github.com/google/uuid"
)

// Concert is a live/concert event discovered for an artist.
type Concert struct {
	ID            uuid.UUID
	ArtistID      uuid.UUID
	ArtistName    string
	Title         string
	VenueName     string
	VenueLocation string
	Date          time.Time
	CoPerformers  []string
	IsFestival    bool
	SourceURL     string
	RawText       string
	DiscoveredAt  time.Time
	CreatedAt     time.Time
}

// IsNew reports whether the concert was discovered within the given window
// (measured from now), per the dashboard's "NEW" badge rule.
func (c Concert) IsNew(now time.Time, window time.Duration) bool {
	return now.Sub(c.DiscoveredAt) <= window
}
