package types

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	PLAYER UserRole = "PLAYER"
	HELPER UserRole = "HELPER"
	ADMIN  UserRole = "ADMIN"
)

type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Members []UserProfile `json:"members"`
}
