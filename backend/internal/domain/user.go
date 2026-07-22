package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is a registered user of the notifier.
type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
}
