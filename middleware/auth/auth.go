package auth

import (
	"fmt"
	"rashomon/api/response"
	"rashomon/model"
	"rashomon/pkg/hash"
	"rashomon/service/user"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var Auth *jwt.GinJWTMiddleware

func init() {
	var err error
	Auth, err = New()
	if err != nil {
		panic(fmt.Sprintf("dashboard auth middle err:%+v", err))
	}
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Id   uint64
	Name string
}

func New() (*jwt.GinJWTMiddleware, error) {
	identityKey := "user"
	authMiddleWare, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "name",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 60,
		MaxRefresh:  time.Minute * 10,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: &User{
						Id:   v.ID,
						Name: v.Name,
					},
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			calims := jwt.ExtractClaims(c)
			userInfo := &User{}
			if userIdentity, ok := calims[identityKey].(*model.User); ok {
				userInfo.Id = userIdentity.ID
				userInfo.Name = userIdentity.Name
			}
			return userInfo
		},
		Authenticator: authenticator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.Error(c, code, message)
		},
		LoginResponse:   loginResponse,
		LogoutResponse:  logoutResponse,
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

func authenticator(c *gin.Context) (interface{}, error) {
	loginVal := login{}
	if err := c.ShouldBind(&loginVal); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userInfo, err := user.NewUserService().GetUserByName(c, loginVal.Username)
	if err != nil {
		return "", err
	}
	if !hash.BcryptCheck(loginVal.Password, userInfo.Password) {
		return nil, jwt.ErrFailedAuthentication
	}
	return userInfo, nil
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	data := map[string]interface{}{
		"token":   token,
		"expires": expire.Format(time.RFC3339),
	}
	response.Success(c, data)
}

func logoutResponse(c *gin.Context, code int) {
	response.Success(c, "success")
}

func GetUserAuthInfo(c *gin.Context) *User {
	return Auth.IdentityHandler(c).(*User)
}
