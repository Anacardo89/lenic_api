package model

import "time"

type Post struct {
	Id        int
	GUID      string
	AuthorId  int
	Title     string
	Content   string
	Image     string
	ImageExt  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsPublic  bool
	Rating    int
	Active    int
}

type PostRatings struct {
	PostId      int
	UserId      int
	RatingValue int
}
