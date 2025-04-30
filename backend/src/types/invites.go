package types

import "time"

type RegistrationInvite struct {
	InviteId  ID        `json:"invite_id" db:"inviteid"`
	Author    ID        `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Expires   time.Time `json:"expires" db:"expires"`
	Code      string    `json:"code" db:"code"`
}
