package backendusers

type CreateBackendUserResponse struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}
