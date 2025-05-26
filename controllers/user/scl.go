package user

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/common"
	"mental/models"
	"mental/service"
	"strconv"
)

// SCLController 处理 SCL 相关请求
type SCLController struct {
	common.BaseController
}

// InsertSCL 处理用户提交 SCL 记录
// @Summary 处理用户提交 SCL 记录，插入scl到数据库
// @Description 插入单个scl记录接口
// @Tags 管理员/用户
// @Produce json
// @Router /scl [post]
func (con SCLController) InsertSCL(c *gin.Context) {
	var scl models.SCL

	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}

	// 绑定 JSON 并检查错误
	if err := c.ShouldBindJSON(&scl); err != nil {
		con.Error(c, nil, "参数格式错误: "+err.Error())
		return
	}

	scl.StudentID = &userId

	// 创建 SCLService 实例
	sclService := service.NewSCLService()

	// 进行业务逻辑校验和存储
	if err := sclService.CreateSCL(&scl); err != nil {
		con.Error(c, nil, err.Error())
		return
	}

	// 返回成功消息
	con.Success(c, "SCL记录创建成功")
}

// SelectAllByUserId 根据用户id查询用户所有历史评测记录
// @Summary 根据当前用户id，查询该用户所有的历史记录数据
// @Description 查询当前用户的历史评测记录列表
// @Tags 管理员/用户
// @Produce json
// @Router /scl [get]
func (con SCLController) SelectAllByUserId(c *gin.Context) {
	id, _ := c.Get("id")
	userId, ok := id.(int64)
	if !ok {
		con.Error(c, nil, "无效的用户id")
		return
	}
	sclService := service.NewSCLService()
	scls, err := sclService.SelectAllByUserId(userId)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, scls)
}

// SelectSCLs 查询所有的scl记录数据
// @Summary 查询所有用户的scl记录数据
// @Description 查询所有用户的历史评测记录列表
// @Tags 管理员
// @Produce json
// @Router /scl/all [get]
func (con SCLController) SelectSCLs(c *gin.Context) {
	sclService := service.NewSCLService()
	scls, err := sclService.SelectAll()
	if err != nil {
		con.Error(c, nil, "查询所有数据记录失败")
		return
	}
	con.Success(c, scls)
}

// DeleteSCL 删除指定的SCL数据
// @Summary 删除指定id的SCL数据
// @Description 删除指定id的SCL数据
// @Tags 管理员
// @Produce json
// @Router /scl [delete]
func (con SCLController) DeleteSCL(c *gin.Context) {
	idStr := c.Query("id") // 获取的是字符串
	sclId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		con.Error(c, nil, "无效的SCL记录id")
		return
	}

	sclService := service.NewSCLService()
	err = sclService.DeleteSCL(sclId)
	if err != nil {
		con.Error(c, nil, err.Error())
		return
	}
	con.Success(c, nil)
}

// UpdateSCL 编辑指定的scl数据
// @Summary 编辑指定的scl数据
// @Description 编辑指定的scl数据
// @Tags 管理员
// @Produce json
// @Router /scl/update [post]
func (con SCLController) UpdateSCL(c *gin.Context) {
	var scl models.SCL
	if err := c.ShouldBindJSON(&scl); err != nil {
		con.Error(c, nil, "参数绑定失败")
		return
	}
	sclService := service.NewSCLService()
	err := sclService.UpdateSCL(&scl)
	if err != nil {
		con.Error(c, nil, "更新数据失败")
		return
	}
	con.Success(c, nil)

}
