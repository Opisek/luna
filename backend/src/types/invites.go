package types

import "time"

type RegistrationInvite struct {
	InviteId  ID        `json:"invite_id" db:"inviteid"`
	Author    ID        `json:"author" db:"author"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Expires   time.Time `json:"expires_at" db:"expires_at"`
	Code      string    `json:"code" db:"code"`
}
