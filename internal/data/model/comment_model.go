package model

import "time"

type Comment struct {
	Id        int
	PostGUID  string
	AuthorId  int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Rating    int
	Active    int
}

type CommentRatings struct {
	CommentId   int
	UserId      int
	RatingValue int
}
