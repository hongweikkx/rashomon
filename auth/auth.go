package auth

import (
	"github.com/hongweikkx/rashomon/hystrix"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"net/http"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func New() *jwt.GinJWTMiddleware{
	identityKey := "id"
	authMiddleWare, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "name",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 10,
		MaxRefresh:  time.Minute * 10,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			calims := jwt.ExtractClaims(c)
			return &User{
				UserName: calims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userId := loginVals.Username
			password := loginVals.Password
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return &User{
					UserName:  userId,
					LastName:  "hongwei",
					FirstName: "gao",
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c * gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code": code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup:  "header: Authorization, query: token, cookie: jwt",
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

	})
	if err != nil {
		log.SugarLogger.Fatal("JWT Error:", err.Error())
	}
	return authMiddleWare
}


func Use(authMiddleware *jwt.GinJWTMiddleware, engine *gin.Engine) {
	// 404
	engine.NoRoute(authMiddleware.MiddlewareFunc())
	// auth group
	authGroup := engine.Group("/auth")
	// Refresh time can be longer than token timeout
	authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/hello", hystrix.HandleDegrade, helloHandler)
	}
	engine.POST("/login", authMiddleware.LoginHandler)
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message" : "hello, world!",
	})
}
