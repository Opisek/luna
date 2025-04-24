package types

import (
	"net"
	"time"
)

type Session struct {
	SessionId    ID        `json:"session_id" db:"sessionid"`
	UserId       ID        `json:"user_id" db:"userid"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastSeen     time.Time `json:"last_seen" db:"last_seen"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IpAddress    net.IP    `json:"ip_address" db:"ip_address"`
	IsShortLived bool      `json:"is_short_lived" db:"is_short_lived"`
	IsApi        bool      `json:"is_api" db:"is_api"`
	SecretHash   []byte    `json:"-" db:"hash"`
}
