package models

type UserRequest struct {
	Username string `json:"username,omitempty" `

	Password string `json:"password,omitempty" `

	Name string `json:"name,omitempty" `
}

func (*UserRequest) CollectionName() string {
	return "user"
}
