package models

import "time"

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

type Post struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name,omitempty"`
	Content    string    `json:"content"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type UpdatePostRequest struct {
	Content    *string `json:"content"`
	Visibility *string `json:"visibility"`
}

type PostResponse struct {
	Post Post `json:"post"`
}

type PostsListResponse struct {
	Posts []Post `json:"posts"`
}
