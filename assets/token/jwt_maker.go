package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey []byte
}

// write new function to create new token using jwt

func NewJWTMaker(secret string) (Maker, error) {

	maker := &JWTMaker{
		secretKey: []byte(secret),
	}
	return maker, nil
}
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (*Payload, string, error) {
	payload := CreateNewPayload(username, duration)
	claims := jwt.MapClaims{
		"sub": payload.Sub,
		"iss": payload.Iss,
		"aud": payload.Aud,
		"exp": payload.Exp,
		"iat": payload.Iat,
	}
	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(maker.secretKey)
	if err != nil {
		return nil, "", err
	}

	return payload, signedToken, nil

}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return maker.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Map claims back to Payload
		// log.Printf("claims: %v", claims)
		payload := &Payload{
			Sub: claims["sub"].(string),
			Iss: claims["iss"].(string),
			Aud: claims["aud"].(string),
			Exp: int64(claims["exp"].(float64)),
			Iat: int64(claims["iat"].(float64)),
		}

		// Check expiration manually (optional, but recommended)
		if !payload.Valid() {
			return nil, fmt.Errorf("token is expired")
		}

		return payload, nil
	}

	return nil, fmt.Errorf("invalid token")

}
