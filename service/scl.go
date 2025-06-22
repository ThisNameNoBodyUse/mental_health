package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"mental/config"
	"mental/dao"
	"mental/models"
	"mental/vo"
	"strings"
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

// CreateSCL 插入 SCL 记录并进行数据校验
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
func (s *SCLService) SelectAllByUserId(userId int64) (*vo.UserSCLResult, error) {
	sclDao := dao.NewSCLDao(config.DB)
	scls, err := sclDao.SelectAllByUserId(userId)
	if err != nil {
		return nil, err
	}

	var (
		results          []vo.SCLRecordAnalysisVO
		recordCount      int
		sumTotalScore    float64
		sumPositiveItems float64
		sumFactors       = make([]float32, 10) // 因子总和
	)

	for _, scl := range scls {
		// 每条记录计算单独健康状态
		healthStatus := calculateHealthStatus(scl)
		results = append(results, vo.SCLRecordAnalysisVO{
			SCL:          scl,
			HealthStatus: healthStatus,
		})

		sumTotalScore += scl.TotalScore
		sumPositiveItems += scl.PositiveItems
		recordCount++

		// 累加每个因子
		sumFactors[0] += scl.Somatization
		sumFactors[1] += scl.Obsession
		sumFactors[2] += scl.Interpersonal
		sumFactors[3] += scl.Depression
		sumFactors[4] += scl.Anxiety
		sumFactors[5] += scl.Hostility
		sumFactors[6] += scl.Phobia
		sumFactors[7] += scl.Paranoia
		sumFactors[8] += scl.Psychoticism
		sumFactors[9] += scl.Other
	}

	// 构造“平均测评”
	var overallHealth string
	if recordCount > 0 {
		avgFactors := make([]float32, 10)
		for i := 0; i < 10; i++ {
			avgFactors[i] = sumFactors[i] / float32(recordCount)
		}
		avgSCL := models.SCL{
			TotalScore:    sumTotalScore / float64(recordCount),
			PositiveItems: sumPositiveItems / float64(recordCount),
			Somatization:  avgFactors[0],
			Obsession:     avgFactors[1],
			Interpersonal: avgFactors[2],
			Depression:    avgFactors[3],
			Anxiety:       avgFactors[4],
			Hostility:     avgFactors[5],
			Phobia:        avgFactors[6],
			Paranoia:      avgFactors[7],
			Psychoticism:  avgFactors[8],
			Other:         avgFactors[9],
		}
		// 计算总体心理状态
		overallHealth = calculateHealthStatus(avgSCL)
	}

	return &vo.UserSCLResult{
		Records:           results,
		UserOverallHealth: overallHealth,
	}, nil
}

// SelectAll 查询所有用户的scl数据
func (sclService *SCLService) SelectAll() ([]models.SCL, error) {
	sclDao := dao.NewSCLDao(config.DB)
	return sclDao.SelectAll()
}

// calculateHealthStatus 综合判断心理健康状态，优先级：总分 > 阳性项目数 > 因子得分
func calculateHealthStatus(scl models.SCL) string {
	// 总分
	totalScore := scl.TotalScore
	// 阳性项目数
	positiveCount := scl.PositiveItems
	// 因子映射
	factorMap := map[string]float32{
		"躯体化":    scl.Somatization,
		"强迫症状":   scl.Obsession,
		"人际关系敏感": scl.Interpersonal,
		"抑郁":     scl.Depression,
		"焦虑":     scl.Anxiety,
		"敌对":     scl.Hostility,
		"恐怖":     scl.Phobia,
		"偏执":     scl.Paranoia,
		"精神病性":   scl.Psychoticism,
		"其他":     scl.Other,
	}
	// 第一优先级：总分判断
	switch {
	case totalScore > 200:
		return "心理状态较差"
	case totalScore > 160:
		return "存在明显心理困扰"
	}
	// 第二优先级：阳性项目数判断
	if positiveCount > 43 {
		return "需关注，阳性项目数较多"
	}
	// 第三优先级：因子均值判断（定位具体问题方向）
	var (
		severeFactors   []string
		moderateFactors []string
		mildFactors     []string
	)
	for name, score := range factorMap {
		switch {
		case score >= 4:
			severeFactors = append(severeFactors, name)
		case score >= 3:
			moderateFactors = append(moderateFactors, name)
		case score >= 2:
			mildFactors = append(mildFactors, name)
		}
	}
	if len(severeFactors) > 0 {
		return "重度" + strings.Join(severeFactors, "、")
	}
	if len(moderateFactors) > 0 {
		return "中度" + strings.Join(moderateFactors, "、")
	}
	if len(mildFactors) > 0 {
		return "轻度" + strings.Join(mildFactors, "、")
	}
	return "心理状态基本正常"
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
