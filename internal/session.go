package internal

import (
	"crypto/rand"
	"encoding/base64"
)

var _ SessionManager = &InMemeorySessiomManagerImpl{}

type SessionManager interface {
	CreateNewSession(value any) (string, error)
	DestroyCurrentSession() error
	GetCurrentSession() (string, error)
}

func NewSessionManager() SessionManager {
	return &InMemeorySessiomManagerImpl{}
}

type InMemeorySessiomManagerImpl struct {
	session string
}

// CreateNewSession implements SessionManager.
func (i *InMemeorySessiomManagerImpl) CreateNewSession(value any) (string, error) {
	session, err := i.randomString(20)
	if err != nil {
		return session, err
	}

	i.session = session
	return session, nil
}

func (i *InMemeorySessiomManagerImpl) randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe, no padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// DestroyCurrentSession implements SessionManager.
func (i *InMemeorySessiomManagerImpl) DestroyCurrentSession() error {
	i.session = ""
	return nil
}

// GetCurrentSession implements SessionManager.
func (i *InMemeorySessiomManagerImpl) GetCurrentSession() (string, error) {
	if i.session == "" {
		return i.CreateNewSession(nil)
	}

	return i.session, nil
}
