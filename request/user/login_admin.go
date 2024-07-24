package user

type (
	LoginRequestAdmin struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" biding:"required"`
	}
	LoginResponseAdmin struct {
		Token string `json:"token"`
	}
)
