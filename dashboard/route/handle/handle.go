package handle

import (
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

func Health(c *gin.Context) {
	c.JSON(200, "")
}

func UserInfo(c *gin.Context) {
	//id := auth.AuthMiddleWare.IdentityHandler(c).(auth.User).Id
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"roles": []string{"admin"},
			"introduction": "I am a super administrator",
			"avatar": "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			"name": "Super Admin",
		},
	})
}

type User struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Status string `json:"status"`
}

func TableList(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"total": 1,
			"items": []User{
				{
					Id:    "1",
					Title: "1",
					Status: "deleted",
				},
			},
		},
	})
}
