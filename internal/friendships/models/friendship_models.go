package models

import "time"

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

type Friendship struct {
	ID          string    `json:"id"`
	RequesterID string    `json:"requester_id"`
	AddresseeID string    `json:"addressee_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FriendRequest struct {
	AddresseeID string `json:"addressee_id"`
}

type FriendshipResponse struct {
	Friendship Friendship `json:"friendship"`
}

type FriendshipsListResponse struct {
	Friendships []Friendship `json:"friendships"`
}

type FriendsListResponse struct {
	Friends []FriendProfile `json:"friends"`
}

type FriendProfile struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}
