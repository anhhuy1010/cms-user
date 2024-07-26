package user

type (
	GetInsertRequest struct {
		Uuid     string `json:"uuid" `
		UserName string `json:"username"`
		IsActive int    `json:"is_active"`
		Role     string `json:"role"`
	}
	InsertResponse struct {
		Uuid string `json:"uuid" `
	}
)
