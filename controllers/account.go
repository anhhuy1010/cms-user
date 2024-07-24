package controllers

import (
	"fmt"
	"net/http"

	"github.com/anhhuy1010/cms-user/helpers/respond"
	"github.com/anhhuy1010/cms-user/helpers/util"
	"github.com/anhhuy1010/cms-user/models"
	request "github.com/anhhuy1010/cms-user/request/account"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountController struct {
}

func (userCtl AccountController) LoginAdmin(c *gin.Context) {
	adminModel := models.Admins{}
	var req request.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"username": req.UserName}
	admin, err := adminModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("user not found"))
		return
	}

	if admin.Password != req.Password {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("password not found"))
		return
	}

	token, err := util.GenerateJWT(admin.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("found"))
		return
	}

	c.JSON(http.StatusOK, respond.Success(request.LoginResponse{Token: token}, "login successfully"))
}
func (userCtl AccountController) LoginCustomer(c *gin.Context) {
	customerModel := models.CustomersSignUp{}

	var req request.LoginRequestUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"username": req.UserName}
	customer, err := customerModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("user not found"))
		return
	}
	if customer.Password != req.Password {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("password not found"))
		return
	}

	token, err := util.GenerateJWT(customer.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("found"))
		return
	}
	customerLogin := models.CustomersLogin{}
	customerLogin.CustomerUuid = customer.Uuid
	customerLogin.Uuid = util.GenerateUUID()
	customerLogin.Token = token

	_, err = customerLogin.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(request.LoginResponse{Token: token}, "login successfully"))
}
func (userCtl AccountController) SignUp(c *gin.Context) {

	var req request.SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	customerData := models.CustomersSignUp{}
	customerData.IsActive = 1
	customerData.Uuid = util.GenerateUUID()
	customerData.Username = req.UserName
	customerData.Name = req.Name
	customerData.Email = req.Email
	customerData.Age = req.Age
	customerData.Password = req.Password
	check := req.CheckPass
	if check != customerData.Password {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Passwords are inconsistent"))
		return
	}
	_, err = customerData.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(customerData.Username, "sign up successfully"))
}
