package simpic

import (
	"crypto/rand"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	createAdminUsername = flag.String("create-admin-username", "", "create a new admin user with this username")
	createAdminPassword = flag.String("create-admin-password", "", "create a new admin user with this password")
)

type User struct {
	Id           int    `db:"user_id"`
	Name         string `db:"user_name"`
	PasswordSalt []byte `db:"user_password_salt"`
	PasswordHash []byte `db:"user_password_hash"`
	SessionKey   []byte `db:"user_session_key"`
	Admin        bool   `db:"user_admin"`
}

type UserManager struct {
	db *Database
}

func NewUserManager(db *Database) *UserManager {
	return &UserManager{db: db}
}

func (u *UserManager) CreateAdmin() {
	if len(*createAdminUsername) > 0 && len(*createAdminPassword) > 0 {
		log.Printf("Creating new admin user '%s' from configuration...\n", *createAdminUsername)
		_, err := u.AddUser(*createAdminUsername, *createAdminPassword, true)
		if err != nil {
			log.Printf("Unable to create user '%s': %v\n", *createAdminUsername, err)
		}
	}
}

func (u *UserManager) AddUser(username, password string, admin bool) (*User, error) {
	var (
		passwordSalt = u.randomBytes(16)
		sessionKey   = u.randomBytes(32)
		passwordHash []byte
		err          error
	)

	if passwordHash, err = bcrypt.GenerateFromPassword(u.salted(password, passwordSalt), 0); err != nil {
		return nil, err
	}

	user := &User{
		Name:         username,
		PasswordSalt: passwordSalt,
		PasswordHash: passwordHash,
		SessionKey:   sessionKey,
		Admin:        admin,
	}

	if err = u.db.AddUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManager) randomBytes(len int) []byte {
	res := make([]byte, len)
	n, err := rand.Read(res)

	if n < len || err != nil {
		panic(fmt.Sprintf("Unable to generate random bytes. Wanted: %d, got: %d, err: %s", len, n, err))
	}

	return res
}

func (u *UserManager) salted(password string, salt []byte) []byte {
	return append([]byte(password), salt...)
}
