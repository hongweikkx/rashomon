package handle

import (
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	//id := auth.AuthMiddleWare.IdentityHandler(c).(auth.User).Id
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"roles":        []string{"admin"},
			"introduction": "I am a super administrator",
			"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			"name":         "Super Admin",
		},
	})
}
