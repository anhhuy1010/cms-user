package controllers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"github.com/anhhuy1010/DATN-cms-customer/grpc"
	pbUsers "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/users"

	"github.com/anhhuy1010/DATN-cms-customer/helpers/respond"
	"github.com/anhhuy1010/DATN-cms-customer/helpers/util"
	"github.com/anhhuy1010/DATN-cms-customer/models"
	request "github.com/anhhuy1010/DATN-cms-customer/request/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) List(c *gin.Context) {

	userModel := new(models.Customer)
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
	userModel := new(models.Customer)
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
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// ////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Update(c *gin.Context) {
	userModel := new(models.Customer)
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
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Delete(c *gin.Context) {
	userModel := new(models.Customer)
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
	userModel := new(models.Customer)
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
	userData := models.Customer{}
	userData.Uuid = util.GenerateUUID()
	userData.Username = req.Username
	userData.Password = req.Password
	userData.Email = req.Email
	userData.IsActive = 1
	userData.Password = req.Password

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
	fmt.Println("🚀 Bắt đầu xử lý login")

	// Bind JSON từ body vào struct
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Lỗi bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	fmt.Printf("📥 Dữ liệu nhận được: %+v\n", req)

	// Tìm user theo email
	condition := bson.M{"email": req.Email}
	userModel := models.Customer{}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Printf("❌ Lỗi tìm user: %v\n", err)
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("user not found"))
		return
	}
	fmt.Printf("✅ Tìm thấy user: %+v\n", user)

	// So sánh mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Printf("❌ Mật khẩu sai: %v\n", err)
		c.JSON(http.StatusUnauthorized, respond.ErrorCommon("invalid password"))
		return
	}
	fmt.Println("🔐 Mật khẩu chính xác")

	// Sinh JWT token
	token, err := util.GenerateJWT(user.Username, nil, nil)
	if err != nil {
		fmt.Printf("❌ Lỗi tạo token: %v\n", err)
		c.JSON(http.StatusInternalServerError, respond.ErrorCommon("failed to generate token"))
		return
	}
	fmt.Println("✅ Token đã được tạo")

	// Lưu token vào DB
	userLogin := models.Tokens{
		UserUuid: user.Uuid,
		Uuid:     util.GenerateUUID(),
		Token:    token,
		IsDelete: 0,
	}

	_, err = userLogin.Insert()
	if err != nil {
		fmt.Printf("❌ Lỗi lưu token: %v\n", err)
		c.JSON(http.StatusInternalServerError, respond.UpdatedFail())
		return
	}
	fmt.Println("✅ Token đã được lưu")

	// Trả về thành công
	c.JSON(http.StatusOK, respond.Success(request.LoginResponse{Token: token}, "login successfully"))
}

////////////////////////////////////////////////////////////////////////////

func (userCtl UserController) SignUp(c *gin.Context) {
	var req request.SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userSignup := models.Customer{}
	userSignup.IsActive = 1
	userSignup.Uuid = util.GenerateUUID()
	userSignup.Username = req.UserName
	userSignup.Password = req.Password
	userSignup.Email = req.Email
	userSignup.StartDay = nil
	userSignup.EndDay = nil

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userSignup.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
		return
	}
	userSignup.Password = string(hashedPassword)

	_, err = userSignup.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}

	c.JSON(http.StatusOK, respond.Success(userSignup.Uuid, "sign up successfully"))
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
	userModel := models.Customer{}

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
		UserUuid:     user.Uuid,
		UserStartDay: *user.StartDay,
		UserEndDay:   *user.EndDay,
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

		c.Set("userStartDay", resp.UserStartDay)
		c.Set("userEndDay", resp.UserEndDay)
		c.Set("userUuid", resp.UserUuid)
		c.Next()
	}
}
