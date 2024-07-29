package user

type (
	SignUpRequestAdmin struct {
		Password string `json:"password" binding:"required"`
		UserName string `json:"username" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
)
