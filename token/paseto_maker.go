package token

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type PasetoMaker struct {
	secretKey string
}

func NewPasetoMaker(secret string) Marker {
	return &PasetoMaker{
		secretKey: secret,
	}
}
func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	token := paseto.NewToken()
	id := uuid.New()
	token.Set("id", id)
	token.SetIssuer("issuer")
	token.SetSubject(username)
	token.SetExpiration(time.Now().Add(duration))
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	key, _ := paseto.V4SymmetricKeyFromBytes([]byte(p.secretKey))
	return token.V4Encrypt(key, nil), nil
}
func (p *PasetoMaker) VerifyToken(encryptedToken string) (*Payload, error) {
	key, _ := paseto.V4SymmetricKeyFromBytes([]byte(p.secretKey))
	parser := paseto.NewParserForValidNow()
	parsedToken, err := parser.ParseV4Local(key, encryptedToken, nil)
	if err != nil {
		fmt.Println("Error decrypting token:", err)
		return nil, err
	}
	var id uuid.UUID
	parsedToken.Get("id", &id)
	sub, _ := parsedToken.GetSubject()
	iat, _ := parsedToken.GetIssuedAt()
	exp, _ := parsedToken.GetExpiration()
	return &Payload{
		ID:  id,
		Sub: sub,
		IAT: iat,
		Exp: exp,
	}, nil

}
