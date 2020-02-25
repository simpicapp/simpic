package internal

import (
	"encoding/hex"
	"flag"
	"log"
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

type SessionManager struct {
	db     *Database
	ticker *time.Ticker
	done   chan bool
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

func NewSessionManager(db *Database) *SessionManager {
	return &SessionManager{
		db:     db,
		ticker: time.NewTicker(time.Hour),
		done:   make(chan bool),
	}
}

func (sm *SessionManager) PeriodicallyPruneInactiveSessions() {
	sm.prune()

	for {
		select {
		case <-sm.ticker.C:
			sm.prune()
		case <-sm.done:
			return
		}
	}
}

func (sm *SessionManager) Stop() {
	sm.ticker.Stop()
	sm.done <- true
}

func (sm *SessionManager) prune() {
	if err := sm.db.DeleteExpiredSessions(); err != nil {
		log.Printf("Unable to delete expired sessions: %v\n", err)
	}
}
