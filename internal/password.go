package internal

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/prometheus/common/log"
	"golang.org/x/crypto/bcrypt"
)

var DefaultPasswordHasher = &Base64BCryptPasswordHasher{}

type PasswordHasher interface {
	Decode(encoded []byte) (decoded []byte, err error)
	Encode(password []byte) (encoded []byte, err error)
	Compare(encoded []byte, password []byte) error
}

type PasswordStore interface {
	Get(username string) ([]byte, bool)
	Set(username, password string) error
	Verify(username, password string) error
}

type PasswordStoreOpts struct {
	hasher PasswordHasher
}

type MemoryPasswordStore struct {
	sync.RWMutex
	passwords map[string][]byte
	hasher    PasswordHasher
}

func NewMemoryPasswordStore(passwords map[string][]byte, opts PasswordStoreOpts) *MemoryPasswordStore {
	var hasher PasswordHasher

	if opts.hasher != nil {
		hasher = opts.hasher
	} else {
		hasher = DefaultPasswordHasher
	}

	return &MemoryPasswordStore{
		passwords: passwords,
		hasher:    hasher,
	}
}

func (store *MemoryPasswordStore) Get(username string) ([]byte, bool) {
	store.RLock()
	defer store.RUnlock()

	hash, ok := store.passwords[username]
	return hash, ok
}

func (store *MemoryPasswordStore) Set(username, password string) error {
	// Not Implemented
	return nil
}

func (store *MemoryPasswordStore) Verify(username, password string) error {
	log.Debugf("looking up: %s", username)
	log.Debugf("%v", store.passwords)
	hash, ok := store.Get(username)
	if !ok {
		log.Debugf("username %s not found", username)
		return fmt.Errorf("account not found: %s", username)
	}

	return store.hasher.Compare(hash, []byte(password))
}

type Base64BCryptPasswordHasher struct{}

func (hasher *Base64BCryptPasswordHasher) Decode(encoded []byte) (decoded []byte, err error) {
	if encoded == nil {
		err = fmt.Errorf("empty password")
		return
	}
	decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	log.Debugf("Decode:")
	log.Debugf("decoded: %v", decoded)
	log.Debugf("encoded: %v", encoded)
	_, err = base64.StdEncoding.Decode(decoded, encoded)
	return
}

func (hasher *Base64BCryptPasswordHasher) Encode(password []byte) (encoded []byte, err error) {
	if password == nil {
		err = fmt.Errorf("empty password")
		return
	}
	bcrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return
	}
	base64.StdEncoding.Encode(encoded, bcrypted)
	return
}

func (hasher *Base64BCryptPasswordHasher) Compare(encoded, password []byte) error {
	log.Debugf("encoded: %s", encoded)
	log.Debugf("password: %s", password)
	decoded, err := hasher.Decode(encoded)
	log.Debugf("decoded: %s", decoded)
	log.Debugf("err: %s", err)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(decoded, []byte(password))
}

// DEPRECATED

func DecodePassword(encoded string) (decoded []byte, err error) {
	if encoded == "" {
		err = fmt.Errorf("empty password")
		return
	}
	decoded, err = base64.StdEncoding.DecodeString(encoded)
	return
}

func ComparePassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
