package models

type Post struct {
	ID           int
	AuthorID     int
	Title        string
	Content      string
	Category     []string
	LikeCount    int
	DislikeCount int
	CommentCount int
	ShortVersion string
}
