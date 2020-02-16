package http

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/csmith/simpic"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	audienceUser  = "simpic.user"
	audienceAdmin = "simpic.admin"
)

var (
	authTokenExpiry   = flag.Duration("auth-token-expiry", time.Hour*24*31, "validity duration of tokens given when a user logs in")
	authKeyFile       = flag.String("auth-key-file", "data/auth.key", "path to the ES256 key to use to sign JWT tokens")
	authKeyAutoCreate = flag.Bool("auth-key-create", true, "whether or not the auth key should be created if it doesn't exist")
)

type claims struct {
	jwt.Claims
	UserId     int    `json:"uid,omitempty"`
	SessionKey []byte `json:"sky,omitempty"`
}

func (s *server) handleAuthenticate() http.HandlerFunc {
	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &LoginData{}
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			log.Printf("Failed to parse JSON body: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := s.db.GetUser(data.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No such user '%s'\n", data.Username)
				w.WriteHeader(http.StatusForbidden)
			} else {
				log.Printf("Unable to retrieve user '%s': %v\n", data.Username, err)
				w.WriteHeader(http.StatusBadRequest)
			}

			return
		}

		if !s.usermanager.CheckPassword(user, data.Password) {
			log.Printf("Bad password for user '%s'\n", data.Username)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token, err := s.generateJWT(user)
		if err != nil {
			log.Printf("Unable to create JWT for user '%s': %v\n", data.Username, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		writeJSON(w, http.StatusOK, LoginResponse{Token: token})
	}
}

func (s *server) createSigner() (jose.Signer, error) {
	key, err := s.loadKey()
	if err != nil {
		return nil, err
	}

	jwk := jose.JSONWebKey{
		Key:       key,
		Use:       "sig",
		Algorithm: "ES256",
		KeyID:     "simpic",
	}

	return jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: jwk}, &jose.SignerOptions{})
}

// generateJWT creates a new JWT for the given user.
func (s *server) generateJWT(user *simpic.User) (string, error) {
	var audience jwt.Audience
	if user.Admin {
		audience = jwt.Audience{audienceUser, audienceAdmin}
	} else {
		audience = jwt.Audience{audienceUser}
	}

	claims := claims{
		Claims: jwt.Claims{
			Audience: audience,
			Subject:  user.Name,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(*authTokenExpiry)),
		},
		SessionKey: user.SessionKey,
		UserId:     user.Id,
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}

// loadKey attempts to load an ECDSA key from disk; if the key can't be read and the auto-create flag is enabled,
// it will attempt to create and write a new key.
func (s *server) loadKey() (*ecdsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(*authKeyFile)
	if err != nil {
		log.Printf("Unable to load existing key from '%s': %v\n", *authKeyFile, err)
		if *authKeyAutoCreate {
			return s.createKey()
		} else {
			return nil, err
		}
	}

	return x509.ParseECPrivateKey(data)
}

// createKey attempts to create a new key using the P256 curve, and write it to disk.
func (s *server) createKey() (*ecdsa.PrivateKey, error) {
	log.Printf("Creating new authentication key...")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	data, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(*authKeyFile, data, os.FileMode(0600))
	return key, err
}
