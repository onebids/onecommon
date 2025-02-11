package model

type UserImageMqDto struct {
	UserID int64
}

type UserImageResultMqDto struct {
	UserID    int64  `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}
