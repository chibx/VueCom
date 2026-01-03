package customer

// JWT Format sent back to the client
type CustomerJWTPayload struct {
	UserId int    `json:"user_id"`
	Fname  string `json:"fname"`
}

type CreateCustomerRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=5"`
	Email       string  `json:"email" validate:"required,email"`
	PhoneNumber *string `json:"phone_number" validate:""`
	Image       *string `json:"image" validate:"http_url"`
	Country     uint    `json:"country" validate:"required"`
	Password    *string `json:"password,omitempty" validate:"required"`
}
