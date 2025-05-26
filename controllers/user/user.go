package user

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/common"
)
import "mental/service"

type UserController struct {
	common.BaseController
}

// Login 登录
// @Summary 登录
// @Description 登录接口
// @Tags 管理员/用户
// @Accept json
// @Produce json
// @Router /login [post]
func (con UserController) Login(c *gin.Context) {
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

// Logout 退出登录 目前只让访问令牌失效
// @Summary 退出登录
// @Description 退出登录接口
// @Tags 管理员/用户
// @Produce json
// @Router /logout [post]
func (con UserController) Logout(c *gin.Context) {
	var userService service.UserService
	// 获取访问令牌和刷新令牌
	token := c.GetHeader("Authorization")
	// 让访问令牌失效
	err := userService.UserLogout(token)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}

// Register 注册
// @Summary 注册
// @Description 注册接口
// @Tags 管理员/用户
// @Accept json
// @Produce json
// @Router /register [post]
func (con UserController) Register(c *gin.Context) {
	var userService service.UserService
	// 绑定表单元素
	if err := c.ShouldBindJSON(&userService); err != nil {
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

// GetUserInfo 获取基本信息
// @Summary 获取基本信息
// @Description 获取基本信息
// @Tags 管理员/用户
// @Produce json
// @Router / [get]
func (con UserController) GetUserInfo(c *gin.Context) {
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
	//avatarURL := utils.GenerateFileURL(c, userInfo.Avatar)
	//// 返回数据时替换头像字段为完整 URL
	//userInfo.Avatar = avatarURL
	con.Success(c, userInfo)
}

// RefreshToken 根据刷新令牌，生成新的访问令牌，并和旧的刷新令牌一起返回
// @Summary 刷新令牌
// @Description 刷新访问令牌
// @Tags 管理员/用户
// @Produce json
// @Router /refresh-token [post]
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

// UpdateUserAvatar 修改头像
// @Summary 修改头像
// @Description 修改头像信息
// @Tags 管理员/用户
// @Produce json
// @Router /avatar [post]
func (con UserController) UpdateUserAvatar(c *gin.Context) {
	var userService service.UserService
	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}
	var form struct {
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		con.Error(c, nil, "参数绑定失败")
		return
	}
	err := userService.UpdateAvatar(userId, form.Avatar)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}

// UpdateUsernameOrEmail 修改用户名/邮箱
// @Summary 修改用户名/邮箱
// @Description 修改用户名/邮箱信息
// @Tags 管理员/用户
// @Produce json
// @Router /base [post]
func (con UserController) UpdateUsernameOrEmail(c *gin.Context) {
	var userService service.UserService
	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}
	var form struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		con.Error(c, nil, "参数绑定失败")
		return
	}
	err := userService.UpdateUsernameOrEmail(userId, form.Username, form.Email)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Description 用户修改密码
// @Tags 管理员/用户
// @Accept json
// @Produce json
// @Router /password [post]
func (con UserController) UpdatePassword(c *gin.Context) {
	var userService service.UserService

	// 获取当前用户ID
	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}

	// 参数绑定
	var form struct {
		OldPassword     string `json:"old_password"`     // 原密码
		NewPassword     string `json:"new_password"`     // 新密码
		ConfirmPassword string `json:"confirm_password"` // 确认密码
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		con.Error(c, nil, "参数绑定失败")
		return
	}

	// 参数校验
	if form.OldPassword == "" || form.NewPassword == "" || form.ConfirmPassword == "" {
		con.Error(c, nil, "所有字段均不能为空")
		return
	}
	if form.OldPassword == form.NewPassword {
		con.Error(c, nil, "新密码不能与原密码相同")
		return
	}
	if form.NewPassword != form.ConfirmPassword {
		con.Error(c, nil, "新密码与确认密码不一致")
		return
	}

	// 调用 service 层修改密码
	err := userService.UpdatePassword(userId, form.OldPassword, form.NewPassword)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}
