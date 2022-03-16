package models

type UserInfo struct {
	Username string `json:"username" `
	Name     string `json:"name"`
	ShopId   string `json:"shopId" `
	Role     string `json:"role"`
}
