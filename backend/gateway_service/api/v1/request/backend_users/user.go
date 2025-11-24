package backendusers

type OnlyID struct {
	ID int `json:"id"`
}

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// JWT Format sent back to the client
type CustomerJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

type sharedUserProps struct {
	FullName        string  `json:"full_name" gorm:"not null;type:varchar(255);index" validate:""`
	UserName        *string `json:"user_name" gorm:"type:varchar(255);index" validate:""`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	PhoneNumber     *string `json:"phone_number" gorm:"type:varchar(20)" validate:""`
	Image           *string `json:"image" validate:"url"`
	Country         uint    `json:"country" gorm:"index"`
	IsEmailVerified bool    `json:"email_verified" gorm:"default:FALSE;not null"`
	Password        *string `json:"password,omitempty" validate:"required"`
}

type ApiCustomer struct {
	sharedUserProps
}

// Base Backend Panel User
type ApiBackendUser struct {
	sharedUserProps
	Role string `json:"role" gorm:"type:varchar(50)"`
}
