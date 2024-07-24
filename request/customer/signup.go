package customer

type (
	SignUpRequestCustomer struct {
		Name      string `json:"name" binding:"required"`
		Age       int    `json:"age" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CheckPass string `json:"checkpass" binding:"required"`
		Email     string `json:"email" binding:"required"`
		UserName  string `json:"username" binding:"required"`
	}
	SignUpResponseCustomer struct {
		UserName string `json:"username"`
	}
)
