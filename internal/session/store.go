package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Name   string
	Key    []byte
	Secure bool
	MaxAge time.Duration
}

type Store struct {
	config Config
}

func NewStore(config Config) *Store {
	if config.Name == "" {
		config.Name = "session"
	}

	return &Store{
		config: config,
	}
}

func (s *Store) Get(req *http.Request, rw http.ResponseWriter) (*Token, error) {
	var token *Token

	cookie, err := req.Cookie(s.config.Name)
	switch {
	// no session cookie, create a new one
	case errors.Is(err, http.ErrNoCookie):
		token, err = NewToken(s.config.Key)
		if err != nil {
			return nil, fmt.Errorf("new token: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("get cookie: %w", err)
	default:
		token, err = ParseToken(cookie.Value, s.config.Key)
		if err != nil {
			return nil, fmt.Errorf("parse token: %w", err)
		}
	}

	return token, nil
}

func (s *Store) Set(rw http.ResponseWriter, token *Token) {
	http.SetCookie(rw, &http.Cookie{
		Name:     s.config.Name,
		Value:    token.String(),
		MaxAge:   int(s.config.MaxAge.Seconds()),
		Secure:   s.config.Secure,
		SameSite: http.SameSiteStrictMode,
	})
}
