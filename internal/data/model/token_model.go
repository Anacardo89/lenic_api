package model

import "time"

type Token struct {
	Id        int
	Token     string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
