package authentication

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
