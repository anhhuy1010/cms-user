package util

import (
	"context"
	"encoding/json"

	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/anhhuy1010/cms-user/config"
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
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {

	cfg := config.GetConfig()
	jwtKeyStr := cfg.GetString("auth.key")
	jwtKey := []byte(jwtKeyStr)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
