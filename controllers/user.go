package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/anhhuy1010/cms-user/constant"
	"github.com/anhhuy1010/cms-user/helpers/respond"
	"github.com/anhhuy1010/cms-user/helpers/util"
	"github.com/anhhuy1010/cms-user/models"
	request "github.com/anhhuy1010/cms-user/request/user"
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
type (
	GetListResponse struct {
		Uuid string `json:"uuid"`
		Task string `json:"task"`
	}
)

func (userCtl UserController) List(c *gin.Context) {
	userModel := new(models.Users)
	var req request.GetListRequest

	// kiểm tra đầu vào
	// gán các tham số truy vấn từ yêu cầu HTTP vào biến reg sử dụng biding.query để chỉ định kiểu
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil { // nếu có lỗi trong quá trình gán, trả về phản hồi http với mã trạng thái lỗi missing params
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{}         // khởi tạo một bản đồ "cond" để chứa các điều kiện truy vấn cho cơ sở dữ liệu
	if req.Username != nil { // nếu trường user name khác rỗng
		cond["username"] = req.Username // lấy username theo username
	}

	if req.IsActive != nil { // tương tự
		cond["is_active"] = req.IsActive
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort) // lấy các tùy chọn phân trang từ yêu cầu, là hàm hỗ trợ lấy các giá trị này
	var respData []request.ListResponse                                                //khởi tạo một slice tên là respData để chứa các phản hồi của danh sách người dùng
	users, err := userModel.Pagination(c, cond, optionsQuery)                          // gọi phương thức Pagination của mô hình users để lấy danh sách người dùng dựa trewen các điều kiện
	for _, user := range users {                                                       //duyệt qua từng người dùng trong danh sach users

		res := request.ListResponse{ // danh sách người dùng được trả về
			Uuid:       user.Uuid,
			ClientUuid: user.ClientUuid,
			Name:       user.Name,
			UserName:   user.Username,
			IsActive:   user.IsActive,
		}
		respData = append(respData, res)
	}
	total, err := userModel.Count(c, cond)
	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

func (userCtl UserController) Detail(c *gin.Context) {
	userModel := new(models.Users)
	var reqUri request.GetDetailUri //khai báo một biến dẫn đến hàm request/user
	// Validation input
	err := c.ShouldBindUri(&reqUri) // hàm dùng để tìm đến đường dẫn uri
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("User no found!"))
		return
	}

	response := request.GetDetailResponse{
		Uuid:     user.Uuid,
		UserName: user.Username,
		Name:     user.Name,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// khởi tạo
func (userCtl UserController) Update(c *gin.Context) {
	userModel := new(models.Users) // tạo một model mới
	var reqUri request.UpdateUri   //tạo biến đưa tới hàm updateuri ở model
	// kiểm tra đầu vào
	err := c.ShouldBindUri(&reqUri) //dùng framwork của gin dẫn đến cái đường dẫn Uri
	if err != nil {                 //câu điều kiện kiểm tra xem việc ràng buộc dữ liệu từ phần đường dẫn có thành công hay không
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.UpdateRequest // câu điều kiện kiểm tra việc ràng buộc dữ liệu file json có thành công hay không
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}  // kiểm tra đường dẫn đến uuid
	user, err := userModel.FindOne(condition) // tìm đến uuid
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("User no found!"))
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
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
		c.JSON(http.StatusOK, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsDelete = constant.DELETE

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "Delete successfully"))
}

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
		c.JSON(http.StatusOK, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsActive = *req.IsActive

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}
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
	userData.Name = req.Name
	_, err = userData.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(userData.Uuid, "update successfully"))
}

// resquest: username, password
// respone: token
func (userCtl UserController) Login(c *gin.Context) {
	userModel := models.Users{}
	// get data request
	var req request.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	// get user from database with username
	condition := bson.M{"username": req.UserName}
	user, err := userModel.FindOne(condition) // get user from database
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("user not found"))
		return
	}
	//check password
	if user.Password != req.Password {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("password not found"))
		return
	}
	//get token with username
	token, err := util.GenerateJWT(user.Username)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("found"))
		return
	}
	//return response
	c.JSON(http.StatusOK, respond.Success(request.LoginResponse{Token: token}, "login successfully"))
}
