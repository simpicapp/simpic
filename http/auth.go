package http

import (
	"database/sql"
	"github.com/simpicapp/simpic"
	"log"
	"net/http"
	"time"
)

const (
	cookieName = "SimpicSession"
)

func (s *server) handleAuthenticate() http.HandlerFunc {
	type LoginData struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &LoginData{}
		if !bind(w, r, data) {
			return
		}

		user, err := s.db.GetUser(data.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No such user '%s'\n", data.Username)
				writeError(w, http.StatusForbidden, "Invalid username/password")
			} else {
				log.Printf("Unable to retrieve user '%s': %v\n", data.Username, err)
				writeError(w, http.StatusInternalServerError, "Unexpected error; please try again")
			}
			return
		}

		if !s.usermanager.CheckPassword(user, data.Password) {
			log.Printf("Bad password for user '%s'\n", data.Username)
			writeError(w, http.StatusForbidden, "Invalid username/password")
			return
		}

		session := simpic.NewSession(user, r.RemoteAddr, r.UserAgent())
		if err := s.db.AddSession(session); err != nil {
			log.Printf("Unable to save token for user '%s': %v\n", data.Username, err)
			writeError(w, http.StatusInternalServerError, "Unexpected error; please try again")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    session.Key,
			Expires:  session.Expires,
			Secure:   false, //TODO: Add a debug flag and toggle this appropriately
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleGetSelf() http.HandlerFunc {
	type GetSelfResponse struct {
		Username string    `json:"username"`
		Admin    bool      `json:"is_admin"`
		Expires  time.Time `json:"session_expires"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxUser).(*simpic.User)
		session := r.Context().Value(ctxSession).(*simpic.Session)
		writeJSON(w, http.StatusOK, GetSelfResponse{
			Username: user.Name,
			Admin:    user.Admin,
			Expires:  session.Expires,
		})
	}
}

func (s *server) clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Expires:  time.Now().Add(time.Hour * -24),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
