package session

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"
)

var (
	ErrInvalidTokenFormat    = errors.New("invalid token format")
	ErrInvalidTokenSignature = errors.New("invalid token signature")
)

type Token struct {
	id        uint64
	signature []byte
}

func NewToken(key []byte) (*Token, error) {
	bid := [8]byte{}
	_, err := rand.Read(bid[:])
	if err != nil {
		return &Token{}, err
	}

	sig := hmac.New(crypto.SHA512.New, key).Sum(bid[:])

	return &Token{
		id:        binary.LittleEndian.Uint64(bid[:]),
		signature: sig,
	}, nil
}

func ParseToken(s string, key []byte) (*Token, error) {
	id, sig, ok := strings.Cut(s, ".")
	if !ok {
		return nil, ErrInvalidTokenFormat
	}

	bid := [8]byte{}
	_, err := base64.RawURLEncoding.Decode(bid[:], []byte(id))
	if err != nil {
		return nil, ErrInvalidTokenFormat
	}

	bsig, err := base64.RawURLEncoding.DecodeString(sig)
	if err != nil {
		return nil, ErrInvalidTokenFormat
	}

	expected := hmac.New(crypto.SHA512.New, key).Sum(bid[:])
	if !hmac.Equal(bsig, expected) {
		return nil, ErrInvalidTokenSignature
	}

	return &Token{
		id:        binary.LittleEndian.Uint64(bid[:]),
		signature: bsig,
	}, nil
}

func (t *Token) ID() uint64 {
	return t.id
}

func (t *Token) String() string {
	bid := [8]byte{}
	binary.LittleEndian.PutUint64(bid[:], t.id)

	var dst []byte
	dst = base64.RawURLEncoding.AppendEncode(dst, bid[:])
	dst = append(dst, '.')
	dst = base64.RawURLEncoding.AppendEncode(dst, t.signature)

	return string(dst)
}
