package backendusers

type RefreshResp struct {
	AccessToken string `json:"access_token"`
}

type CreateBackendUserResponse struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}
