package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) Marker {
	return &JWTMaker{
		secretKey: secretKey,
	}
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	id := uuid.New()
	claims := jwt.MapClaims{
		"id":  id,
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, _ := uuid.Parse(claims["id"].(string))
		iatFloat, _ := claims["iat"].(float64)
		expFloat, _ := claims["exp"].(float64)
		payload := &Payload{
			ID:  id,
			Sub: claims["sub"].(string),
			IAT: time.Unix(int64(iatFloat), 0),
			Exp: time.Unix(int64(expFloat), 0),
		}
		return payload, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
