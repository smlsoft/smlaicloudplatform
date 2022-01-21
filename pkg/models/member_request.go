package models

type MemberRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (*MemberRequest) CollectionName() string {
	return "member"
}

type MemberRequestEdit struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

func (*MemberRequestEdit) CollectionName() string {
	return "member"
}

type MemberRequestPassword struct {
	Password string `json:"password" validate:"required,gte=6"`
}

func (*MemberRequestPassword) CollectionName() string {
	return "member"
}
