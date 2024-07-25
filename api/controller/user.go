package controller

import (
	"github.com/gin-gonic/gin"
	"rashomon/api/response"
	"rashomon/middleware/auth"
	userService "rashomon/service/user"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (userC UserController) UserInfo(c *gin.Context) {
	user := auth.GetUserAuthInfo(c)
	userInfo, err := userService.NewUserService().GetUserByName(c, user.Name)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	response.Success(c, userInfo)
}
