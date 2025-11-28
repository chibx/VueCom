package backendusers

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

// type sharedUserProps struct {
// 	FullName        string  `json:"full_name" gorm:"not null;type:varchar(255);index" validate:""`
// 	UserName        *string `json:"user_name" gorm:"type:varchar(255);index" validate:""`
// 	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
// 	PhoneNumber     *string `json:"phone_number" gorm:"type:varchar(20)" validate:""`
// 	Image           *string `json:"image" validate:"url"`
// 	Country         uint    `json:"country" gorm:"index"`
// 	IsEmailVerified bool    `json:"email_verified" gorm:"default:FALSE;not null"`
// 	Password        *string `json:"password,omitempty" validate:"required"`
// }

// Base Backend Panel User
type CreateBackendUserRequest struct {
	FullName        string  `json:"full_name" validate:"required,min=5"`
	UserName        *string `json:"user_name" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	PhoneNumber     *string `json:"phone_number" validate:""`
	Image           *string `json:"image" validate:"http_url"`
	Country         *uint   `json:"country" validate:"required_if=Role owner"`
	IsEmailVerified bool    `json:"email_verified"`
	Password        *string `json:"password,omitempty" validate:"required"`
	Role            string  `json:"role" validate:"required"`
}
