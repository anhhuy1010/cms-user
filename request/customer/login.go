package customer

type (
	LoginRequestCustomer struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" biding:"required"`
	}
	LoginResponseCustomer struct {
		Token string `json:"token"`
	}
)
