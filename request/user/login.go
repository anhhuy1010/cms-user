package user

type (
	LoginRequest struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" biding:"required"`
	}
	LoginResponse struct {
		Token string `json:"token"`
	}
)
