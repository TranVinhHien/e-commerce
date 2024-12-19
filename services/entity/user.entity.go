package services

import "time"

type Users struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	FullName string    `json:"full_name"`
	IsActive bool      `json:"is_active"`
	CreateAt time.Time `json:"create_at"`
}
