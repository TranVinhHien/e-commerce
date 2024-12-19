package token

import (
	"errors"
	"time"
)

var (
	ErrExpireToken  = errors.New("token has Expired")
	ErrInvalidToken = errors.New("token has Invalid")
)

type Payload struct {
	Sub string `json:"sub"`
	Iss string `json:"iss"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

func CreateNewPayload(username string, duration time.Duration) *Payload {
	payload := Payload{
		Sub: username,
		Iss: "new.project.com",
		Aud: "sau này truyền role vào đây",
		Exp: time.Now().Add(duration).Unix(),
		Iat: time.Now().Unix(),
	}
	return &payload
}

func (payload *Payload) Valid() bool {
	currentTime := time.Now().Unix()
	return currentTime <= payload.Exp
}
