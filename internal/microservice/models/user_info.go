package models

type UserInfo struct {
	Username string `json:"username" `
	Name     string `json:"name"`
	ShopID   string `json:"shopID" `
	Role     string `json:"role"`
}
