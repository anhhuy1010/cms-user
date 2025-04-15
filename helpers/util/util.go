package util

import (
	"context"
	"encoding/json"

	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUUID() (s string) {
	uuidNew, _ := uuid.NewUUID()
	return uuidNew.String()
}

func ShoudBindHeader(c *gin.Context) bool {
	platform := c.Request.Header.Get("X-PLATFORM")
	lang := c.Request.Header.Get("X-LANG")

	if platform == "" || lang == "" {
		return false
	}

	return true
}

func GetNowUTC() time.Time {
	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	return currentTime
}
func DebugJson(value interface{}) {
	fmt.Println(reflect.TypeOf(value).String())
	prettyJSON, _ := json.MarshalIndent(value, "", "    ")
	fmt.Printf("%s\n", string(prettyJSON))
}

func GetKeyFromContext(ctx context.Context, key string) (interface{}, bool) {
	if v := ctx.Value(key); v != nil {
		return v, true
	}

	return nil, false
}

func LogPrint(jsonData interface{}) {
	prettyJSON, _ := json.MarshalIndent(jsonData, "", "")
	fmt.Printf("%s\n", strings.ReplaceAll(string(prettyJSON), "\n", ""))
}

type Claims struct {
	Uuid     string     `json:"uuid"`
	StartDay *time.Time `json:"startday"`
	EndDay   *time.Time `json:"endday"`
	jwt.RegisteredClaims
}

// Valid method để implements jwt.Claims
func (c *Claims) Valid() error {
	// Nếu StartDay không nil và EndDay không nil, kiểm tra xem StartDay có lớn hơn EndDay không
	if c.StartDay != nil && c.EndDay != nil && c.StartDay.After(*c.EndDay) {
		return fmt.Errorf("invalid user date range: startday cannot be after endday")
	}
	return nil
}

// GenerateJWT với StartDay và EndDay có thể nil
func GenerateJWT(uuid string, startday, endday *time.Time) (string, error) {
	// Nếu startday là nil, gán nó là thời gian hiện tại
	if startday == nil {
		now := time.Now()
		startday = &now
	}

	// Nếu endday là nil, gán nó là 24 giờ sau startday
	if endday == nil {
		endday = startday
		*endday = startday.Add(720 * time.Hour)
	}

	cfg := config.GetConfig()
	jwtKeyStr := cfg.GetString("auth.key")
	jwtKey := []byte(jwtKeyStr)

	claims := &Claims{
		Uuid:     uuid,
		StartDay: startday,
		EndDay:   endday,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)}, // Đặt thời gian hết hạn là 24 giờ sau
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},                     // Đặt thời gian phát hành là hiện tại
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
