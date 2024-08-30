package token

import (
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID  uuid.UUID `json:"id"`
	Sub string    `json:"sub"`
	IAT time.Time `json:"iat"`
	Exp time.Time `json:"exp"`
}
