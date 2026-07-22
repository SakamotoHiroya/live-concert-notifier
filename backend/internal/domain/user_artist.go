package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserArtist represents a user's follow relationship with an artist.
type UserArtist struct {
	UserID     uuid.UUID
	ArtistID   uuid.UUID
	FollowedAt time.Time
}
