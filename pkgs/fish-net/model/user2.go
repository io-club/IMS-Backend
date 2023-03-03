package model

import (
	"crypto/rand"
	"demo/glb"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

// User represents the user model
type User2 struct {
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential `gorm:"-"`
}

// NewUser creates and returns a new User
func NewUser2(name string, displayName string) *User2 {

	user := &User2{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName
	// user.credentials = []webauthn.Credential{}

	return user
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u User2) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.id))
	return buf
}

// WebAuthnName returns the user's username
func (u User2) WebAuthnName() string {
	return u.name
}

// WebAuthnDisplayName returns the user's display name
func (u User2) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon is not (yet) implemented
func (u User2) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User2) AddWebAuthnCredential(cred webauthn.Credential) {
	glb.LOG.Info("[cred]: " + string(cred.ID))
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User2) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

type userdb struct {
	users map[string]*User2
	mu    sync.RWMutex
}

var db *userdb

func encodeUserID(id uint64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, id)
	urlBase := base64.StdEncoding.EncodeToString(buf)
	urlBase = strings.ReplaceAll(urlBase, "+", "-")
	urlBase = strings.ReplaceAll(urlBase, "/", "_")
	return urlBase
}

// DB returns a userdb singleton
func UserDB() *userdb {

	if db == nil {
		db = &userdb{
			users: make(map[string]*User2),
		}
	}

	return db
}

// TestUser returns true if the user exists
func (db *userdb) TestUser(name string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, ok := db.users[name]
	return ok
}

func (db *userdb) CreateUser(username string) (*User2, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	user := &User2{
		id:          randomUint64(),
		name:        username,
		displayName: username,
	}

	// if UserDB().TestUser(user.name) {
	// 	return nil, fmt.Errorf("User '%s': exists", user.name)
	// }

	glb.LOG.Info("\n[user]: " + encodeUserID(user.id) + "\n")

	db.users[encodeUserID(user.id)] = user

	// userGet, err := UserDB().GetUser(user.id)

	return user, nil
}

// GetUser returns a *User by the user's username
func (db *userdb) GetUser(baseID string) (*User2, error) {

	glb.LOG.Info("\n[user get]: " + baseID + "\n")
	db.mu.Lock()
	defer db.mu.Unlock()
	user, ok := db.users[baseID]
	if !ok {
		return nil, fmt.Errorf("user does not exist")
	}

	return user, nil
}

// PutUser stores a new user by the user's username
func (db *userdb) PutUser(user *User2) {

	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[user.name] = user
}
