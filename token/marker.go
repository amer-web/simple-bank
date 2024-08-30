package token

import (
	"github.com/amer-web/simple-bank/config"
	"time"
)

type Marker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

func NewMakerToken() Marker {
	key := config.Source.TOKENKEY
	switch config.Source.TOKENDRIVER {
	case "jwt":
		return NewJwtMaker(key)
	case "paseto":
		return NewPasetoMaker(key)
	default:
		return NewJwtMaker(key)
	}
}
