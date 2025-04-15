package user

import "time"

type (
	CheckRoleRequest struct {
		Token string `json:"token"`
	}
	CheckRoleResponse struct {
		UserUuid     string    `json:"user_uuid" `
		UserStartDay time.Time `json:"startday" `
		UserEndDay   time.Time `json:"endday" `
	}
)
