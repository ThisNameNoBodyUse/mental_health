package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"mental/config"
	"mental/dao"
	"mental/models"
)

// SCLService 结构体
type SCLService struct {
	Validator *validator.Validate
}

// NewSCLService 创建新的 SCLService
func NewSCLService() *SCLService {
	return &SCLService{
		Validator: validator.New(),
	}
}

// CreateSCL 插入 SCL 记录并进行数据校验 TODO
func (sclService *SCLService) CreateSCL(scl *models.SCL) error {
	// **手动检查某些必须字段**
	if scl.Name == "" {
		return errors.New("姓名不能为空")
	}
	if scl.Age <= 0 {
		return errors.New("年龄必须为正整数")
	}
	if scl.Gender != 0 && scl.Gender != 1 {
		return errors.New("性别只能是 0（女）或 1（男）")
	}

	// *校验评分字段范围
	fields := []struct {
		value float32
		name  string
	}{
		{scl.Somatization, "躯体化"},
		{scl.Obsession, "强迫症状"},
		{scl.Interpersonal, "人际关系敏感"},
		{scl.Depression, "抑郁"},
		{scl.Anxiety, "焦虑"},
		{scl.Hostility, "敌对"},
		{scl.Phobia, "恐怖"},
		{scl.Paranoia, "偏执"},
		{scl.Psychoticism, "精神病性"},
		{scl.Other, "其他"},
	}

	for _, field := range fields {
		if field.value < 0 || field.value > 5 {
			return errors.New(field.name + " 分数必须在 0.0 - 5.0 之间")
		}
	}

	// 使用 Validator 自动校验
	err := sclService.Validator.Struct(scl)
	if err != nil {
		return err
	}

	// 存入数据库
	sclDao := dao.NewSCLDao(config.DB)
	if err := sclDao.Save(scl); err != nil {
		return err
	}

	return nil
}

// SelectAllByUserId 根据用户id查询历史所有的评测记录返回
func (s *SCLService) SelectAllByUserId(userId int64) ([]models.SCL, error) {
	sclDao := dao.NewSCLDao(config.DB)
	return sclDao.SelectAllByUserId(userId)
}

// SelectAll 查询所有用户的scl数据
func (sclService *SCLService) SelectAll() ([]models.SCL, error) {
	sclDao := dao.NewSCLDao(config.DB)
	return sclDao.SelectAll()
}

// DeleteSCL 删除指定SCL记录
func (sclService *SCLService) DeleteSCL(id int64) error {
	sclDao := dao.NewSCLDao(config.DB)
	return sclDao.DeleteByID(id)
}

// UpdateSCL 更新指定SCL记录
func (s *SCLService) UpdateSCL(scl *models.SCL) error {
	dao := dao.NewSCLDao(config.DB)
	return dao.UpdateByID(scl.ID, scl)
}
