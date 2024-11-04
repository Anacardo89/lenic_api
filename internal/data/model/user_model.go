package model

import (
	"time"
)

type User struct {
	Id            int
	UserName      string
	Email         string
	HashPass      string
	ProfilePic    string
	ProfilePicExt string
	Followers     int
	Following     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Active        int
}

type Follows struct {
	FollowerId int
	FollowedId int
	Status     int
}
