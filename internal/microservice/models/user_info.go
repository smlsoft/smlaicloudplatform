package models

type UserInfo struct {
	Username   string `json:"username" `
	Name       string `json:"name"`
	MerchantId string `json:"merchantId" `
	Role       string `json:"role"`
}
