package admin

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/common"
	"mental/utils"
)
import "mental/service"

type UserController struct {
	common.BaseController
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员登录接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Router /admin/login [post]
func (con UserController) AdminLogin(c *gin.Context) {
	var userService service.UserService
	// 绑定表单到userService中
	if err := c.ShouldBind(&userService); err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	data, err := userService.UserLogin()
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, data) // 登录成功，响应数据
}

// AdminLogout 管理员退出登录 目前只让访问令牌失效
// @Summary 管理员退出登录
// @Description 管理员退出登录接口
// @Tags 管理员
// @Produce json
// @Router /admin/logout [post]
func (con UserController) AdminLogout(c *gin.Context) {
	var userService service.UserService
	// 获取访问令牌和刷新令牌
	token := c.GetHeader("Authorization")
	// 让访问令牌失效
	err := userService.AdminLogout(token)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}

// AdminRegister 管理员注册
// @Summary 管理员注册
// @Description 管理员注册接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Router /admin/register [post]
func (con UserController) AdminRegister(c *gin.Context) {
	var userService service.UserService
	// 绑定表单元素
	if err := c.ShouldBind(&userService); err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	res, err := userService.UserRegister()
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, res)

}

// GetAdminInfo 获取管理员基本信息
// @Summary 获取管理员信息
// @Description 获取管理员基本信息
// @Tags 管理员
// @Produce json
// @Router /admin [get]
func (con UserController) GetAdminInfo(c *gin.Context) {
	var userService service.UserService
	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}
	userInfo, err := userService.GetUserInfoById(userId)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	// 拼接头像完整 URL
	avatarURL := utils.GenerateFileURL(c, userInfo.Avatar)
	// 返回数据时替换头像字段为完整 URL
	userInfo.Avatar = avatarURL
	con.Success(c, userInfo)
}

// RefreshToken 根据刷新令牌，生成新的访问令牌，并和旧的刷新令牌一起返回
// @Summary 刷新令牌
// @Description 刷新管理员访问令牌
// @Tags 管理员
// @Produce json
// @Router /admin/refresh-token [post]
func (con UserController) RefreshToken(c *gin.Context) {
	var userService service.UserService
	refreshToken := c.GetHeader("Refresh-Token") // 刷新令牌
	if refreshToken == "" {
		con.Error(c, nil, "Refresh-Token header is required")
		return
	}
	newAccessToken, err := userService.RefreshToken(refreshToken)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, newAccessToken)
}
