package user

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/common"
	"mental/models"
	"mental/service"
)

// SCLController 处理 SCL 相关请求
type SCLController struct {
	common.BaseController
}

// InsertSCL 处理用户提交 SCL 记录
// @Summary 处理用户提交 SCL 记录，插入scl到数据库
// @Description 插入单个scl记录接口
// @Tags 提交 SCL记录
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
// @Tags 查询scl记录
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
