package simpic

import (
	"encoding/hex"
	"flag"
	"time"
)

type Session struct {
	Key       string    `db:"session_key"`
	Created   time.Time `db:"session_created"`
	Expires   time.Time `db:"session_expires"`
	Ip        string    `db:"session_ip"`
	UserAgent string    `db:"session_user_agent"`
	UserId    int       `db:"user_id"`
}

type SessionUser struct {
	Session
	User
}

var (
	sessionExpiry = flag.Duration("session-expiry", time.Hour*24*31, "length of time users stay logged in")
)

func NewSession(user *User, ip, userAgent string) *Session {
	return &Session{
		Key:       hex.EncodeToString(randomBytes(64)),
		Created:   time.Now(),
		Expires:   time.Now().Add(*sessionExpiry),
		Ip:        ip,
		UserAgent: userAgent,
		UserId:    user.Id,
	}
}
