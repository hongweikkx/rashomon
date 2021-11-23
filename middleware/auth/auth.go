package auth

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var AuthMiddleWare *jwt.GinJWTMiddleware

func init() {
	var err error
	AuthMiddleWare, err = New()
	if err != nil {
		panic(fmt.Sprintf("dashboard auth middle err:%+v", err))
	}
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Id string
}

func New() (*jwt.GinJWTMiddleware, error) {
	identityKey := "id"
	authMiddleWare, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "name",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 60,
		MaxRefresh:  time.Minute * 10,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			calims := jwt.ExtractClaims(c)
			return &User{
				Id: calims[identityKey].(string),
			}
		},
		Authenticator: authenticator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"code":    40100,
				"message": message,
			})
		},
		LoginResponse: loginResponse,
		LogoutResponse: logoutResponse,
		RefreshResponse: loginResponse,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		return nil, err
	}
	return authMiddleWare, nil
}


func authenticator(c *gin.Context) (interface{}, error){
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userId := loginVals.Username
	password := loginVals.Password
	if userId == "admin" && password == "123456" {
		return &User{
			Id:  "100000",
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"token": token,
			"expires": expire.Format(time.RFC3339),
		},
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": "success",
	})
}