package user

type (
	SignUpRequestAdmin struct {
		Password string `json:"password" binding:"required"`
		UserName string `json:"username" binding:"required"`
	}
	SignUpResponseAdmin struct {
		UserName string `json:"username"`
	}
)
