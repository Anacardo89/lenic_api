package model

import "time"

type Session struct {
	Id        int
	SessionId string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    int
}
