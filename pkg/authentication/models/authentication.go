package models

import "time"

type AuthenticationContext struct {
	Ip string
}

type ShopFavoriteRequest struct {
	ShopID     string `json:"shopid" bson:"shopid"`
	IsFavorite bool   `json:"isfavorite" bson:"isfavorite"`
}

type TokenLoginRequest struct {
	Token string `json:"token" validate:"required"`
}

type TokenLoginResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

type PhoneNumberLoginReponse struct {
	RefCode string    `json:"refcode"`
	Expire  time.Time `json:"expire"`
}

type PhoneNumberLoginRequest struct {
	PhoneNumber string `json:"phonenumber" bson:"phonenumber" validate:"required,max=233"`
}

type PhoneNumberOTPRequest struct {
	PhoneNumber string `json:"phonenumber" bson:"phonenumber" validate:"required,max=233"`
	RefCode     string `json:"refcode"`
	OTP         string `json:"otp" bson:"otp" validate:"required,max=20"`
}

type PhoneOTP struct {
	PhoneNumber string `json:"phonenumber"`
	OTP         string `json:"otp" `
}
