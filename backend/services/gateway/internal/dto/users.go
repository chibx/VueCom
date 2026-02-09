package dto

type UserForLogin struct {
	ID           uint
	UserName     string
	PasswordHash string
	CreatedBy    *uint
	Role         string
}
