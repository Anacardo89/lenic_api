package model

import "time"

type Notification struct {
	Id         int
	UserID     int
	FromUserId int
	NotifType  string
	NotifMsg   string
	ResourceId string
	ParentId   string
	IsRead     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
