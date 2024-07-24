package controllers

import (
	"fmt"
	"net/http"

	"github.com/anhhuy1010/cms-user/helpers/respond"
	"github.com/anhhuy1010/cms-user/helpers/util"
	"github.com/anhhuy1010/cms-user/models"

	request "github.com/anhhuy1010/cms-user/request/customer"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CustomerController struct {
}

func (customerCtl CustomerController) LoginCustomer(c *gin.Context) {
	customerModel := models.CustomersSignUp{}

	var req request.LoginRequestCustomer
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
		c.JSON(http.StatusOK, respond.ErrorCommon("customer not found"))
		return
	}
	cond := bson.M{"password": req.Password}
	_, err = customerModel.FindOne(cond)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("customer not found"))
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
	c.JSON(http.StatusOK, respond.Success(request.LoginResponseCustomer{Token: token}, "login successfully"))
}
func (customerCtl CustomerController) SignUpCustomer(c *gin.Context) {

	var req request.SignUpRequestCustomer
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	customerSignup := models.CustomersSignUp{}
	customerSignup.IsActive = 1
	customerSignup.Uuid = util.GenerateUUID()
	customerSignup.Username = req.UserName
	customerSignup.Name = req.Name
	customerSignup.Email = req.Email
	customerSignup.Age = req.Age
	customerSignup.Password = req.Password
	check := req.CheckPass
	if check != customerSignup.Password {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Passwords are inconsistent"))
		return
	}
	_, err = customerSignup.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(customerSignup.Username, "sign up successfully"))
}
