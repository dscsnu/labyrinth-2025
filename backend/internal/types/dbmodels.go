package types

import (
	"github.com/google/uuid"
	"time"
)

type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
