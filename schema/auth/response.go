package auth

import "time"

type AuthRes struct {
	Username              string    `json:"username"`
	AccessToken           *string   `json:"access_token"`
	RefreshToken          *string   `json:"refresh_token"`
	Fullname              string    `json:"fullname"`
	Iat                   time.Time `json:"iat"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}
