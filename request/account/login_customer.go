package account

type (
	LoginRequestUser struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" biding:"required"`
	}
	LoginResponseUser struct {
		Token string `json:"token"`
	}
)
