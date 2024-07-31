package controllers

import (
	"context"
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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type UserController struct {
}

// List
// @Summary Get list users test ss
// @Schemes
// @Description Get list users
// @Tags users
// @Accept json
// @Produce json
// @Param request query request.GetListRequest true "query params"
// @Success 200 {object} respond.PaginationResponse
// @Router /users [get]

// khởi tạo
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
		UserName: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// //////////////////////////////////////////////////////////////////////////
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

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
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

	if *req.IsActive < 0 || *req.IsActive >= 5 {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Stauts just can be set in range [0..5]"))
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
	userData.Username = req.UserName
	userData.Uuid = req.Uuid
	_, err = userData.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(userData.Uuid, "update successfully"))
}

// ////////////////////////////////////////////////////
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
func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		userCtl := UserController{}
		resp, err := userCtl.CheckRole(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userRole", resp.Role)
		c.Set("userUuid", resp.UserUuid)
		if resp.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}
