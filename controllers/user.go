package controllers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/anhhuy1010/cms-user/constant"
	"github.com/anhhuy1010/cms-user/grpc"

	pbUsers "github.com/anhhuy1010/cms-user/grpc/proto/users"
	"github.com/anhhuy1010/cms-user/helpers/respond"
	"github.com/anhhuy1010/cms-user/helpers/util"
	"github.com/anhhuy1010/cms-user/models"
	request "github.com/anhhuy1010/cms-user/request/user"
	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type UserController struct {
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) List(c *gin.Context) {
	userModel := new(models.Users)
	var req request.GetListRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{}
	if req.Username != nil {
		cond["username"] = req.Username
	}

	if req.IsActive != nil {
		cond["is_active"] = req.IsActive
	}
	if req.Role != nil {
		cond["role"] = req.Role
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)
	var respData []request.ListResponse
	users, err := userModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	for _, user := range users {
		res := request.ListResponse{
			Uuid:     user.Uuid,
			UserName: user.Username,
			IsActive: user.IsActive,
			Role:     user.Role,
		}
		respData = append(respData, res)
	}
	total, err := userModel.Count(c, cond)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Detail(c *gin.Context) {
	userModel := new(models.Users)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	response := request.GetDetailResponse{
		Uuid:     user.Uuid,
		Username: user.Username,
		Email:    user.Email,
		IsActive: user.IsActive,
		IsDelete: user.IsDelete,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// ////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Update(c *gin.Context) {
	userModel := new(models.Users)
	var reqUri request.UpdateUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.UpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.UserName != "" {
		user.Username = req.UserName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Delete(c *gin.Context) {
	userModel := new(models.Users)
	var reqUri request.DeleteUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsDelete = constant.DELETE

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "Delete successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) UpdateStatus(c *gin.Context) {
	userModel := new(models.Users)
	var reqUri request.UpdateStatusUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.UpdateStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if *req.IsActive < 0 || *req.IsActive > 1 {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Stauts just can be set in range [0..1]"))
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsActive = *req.IsActive

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Create(c *gin.Context) {
	var req request.GetInsertRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userData := models.Users{}
	userData.Uuid = util.GenerateUUID()
	userData.Username = req.Username
	userData.Password = req.Password
	userData.Email = req.Email
	userData.IsActive = 1
	userData.Password = req.Password
	userData.Role = "user"
	userData.CreatedAt = time.Time{}
	userData.UpdatedAt = time.Time{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
		return
	}
	userData.Password = string(hashedPassword)

	_, err = userData.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(userData.Uuid, "create successfully"))
}

// /////////////////////////////////////////////////////
func (userCtl UserController) Login(c *gin.Context) {
	adminModel := models.Users{}

	var req request.LoginRequestAdmin
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
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("user not found"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
		return
	}
	token, err := util.GenerateJWT(admin.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("found"))
		return
	}
	adminLogin := models.Tokens{}
	adminLogin.UserUuid = admin.Uuid
	adminLogin.Uuid = util.GenerateUUID()
	adminLogin.Token = token
	adminLogin.IsDelete = 0

	_, err = adminLogin.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(request.LoginResponseAdmin{Token: token}, "login successfully"))
}

////////////////////////////////////////////////////////////////////////////

func (userCtl UserController) SignUp(c *gin.Context) {
	var req request.SignUpRequestAdmin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	adminSignup := models.Users{}
	adminSignup.IsActive = 1
	adminSignup.Uuid = util.GenerateUUID()
	adminSignup.Username = req.UserName
	adminSignup.Password = req.Password
	adminSignup.Role = req.Role
	adminSignup.Email = req.Email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminSignup.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
		return
	}
	adminSignup.Password = string(hashedPassword)

	_, err = adminSignup.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}

	c.JSON(http.StatusOK, respond.Success(adminSignup.Uuid, "sign up successfully"))
}

func (userCtl UserController) CheckRole(token string) (*pbUsers.DetailResponse, error) {
	grpcConn := grpc.GetInstance()
	client := pbUsers.NewUserClient(grpcConn.UsersConnect)
	req := pbUsers.DetailRequest{
		Token: token,
	}
	resp, err := client.Detail(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (userCtl UserController) GetRoleByToken(token string) (*request.CheckRoleResponse, error) {
	tokenModel := models.Tokens{}
	userModel := models.Users{}

	condition := bson.M{"token": token}
	tokenDoc, err := tokenModel.FindOne(condition)
	if err != nil {
		return nil, errors.New("token not found")
	}
	if tokenDoc == nil {
		return nil, errors.New("token document is nil")
	}

	cond := bson.M{"uuid": tokenDoc.UserUuid}
	user, err := userModel.FindOne(cond)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := &request.CheckRoleResponse{
		UserUuid: user.Uuid,
		Role:     user.Role,
	}
	return resp, nil
}

func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		userCtl := UserController{}
		resp, err := userCtl.GetRoleByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userRole", resp.Role)
		c.Set("userUuid", resp.UserUuid)

		if resp.Role != "admin" && c.Request.Method != http.MethodGet {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
