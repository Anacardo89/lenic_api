package model

type Tag struct {
	Id      int
	TagName string
	TagType string
}

type UserTag struct {
	Id        int
	TagId     int
	PostId    int
	CommentId int
	TagPlace  string
}

type ReferenceTag struct {
	Id        int
	TagId     int
	PostId    int
	CommentId int
	TagPlace  string
}
