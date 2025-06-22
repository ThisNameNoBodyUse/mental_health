package vo

import "mental/models"

// 单次测评分析结构体，嵌入SCL，增加计算字段
type SCLRecordAnalysisVO struct {
	models.SCL          // SCL基础结构体
	HealthStatus string `json:"health_status"` // 测评整体状态
}

type UserSCLResult struct {
	Records           []SCLRecordAnalysisVO `json:"records"`       // 用户每次测评记录
	UserOverallHealth string                `json:"health_result"` // 整体心理状态
}
