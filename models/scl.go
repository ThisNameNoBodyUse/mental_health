package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// CustomTime 让 TestDate 解析为 YYYY-MM-DD
type CustomTime time.Time

const customFormat = "2006-01-02"

// UnmarshalJSON 解析 JSON -> time.Time
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1] // 去掉 JSON 的 `"`
	parsedTime, err := time.Parse(customFormat, str)
	if err != nil {
		return err
	}
	*t = CustomTime(parsedTime)
	return nil
}

// MarshalJSON 解析 time.Time -> JSON
func (t CustomTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(customFormat)
	return json.Marshal(formatted)
}

// 实现 driver.Value 接口，让 GORM 识别 CustomTime
func (t CustomTime) Value() (driver.Value, error) {
	return time.Time(t), nil // 让 GORM 存储 time.Time 类型
}

// 实现 Scan 方法，让 GORM 从数据库解析 CustomTime
func (t *CustomTime) Scan(value interface{}) error {
	if v, ok := value.(time.Time); ok {
		*t = CustomTime(v)
		return nil
	}
	return nil
}

// SCL 表示 scl 表的结构体，记录 SCL-90 心理测评记录
type SCL struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	StudentID *int64     `json:"student_id,omitempty" gorm:"column:student_id;comment:学生ID"` // 可空
	Name      string     `json:"name" gorm:"type:varchar(50);not null;comment:学生姓名"`
	Gender    int        `json:"gender" gorm:"type:tinyint;not null;comment:性别 0女 1男"`
	Age       int        `json:"age" gorm:"not null;comment:年龄"`
	TestDate  CustomTime `json:"test_date" gorm:"type:date;not null;comment:测评日期"`

	Somatization  float32 `json:"somatization" gorm:"type:decimal(3,1);not null;comment:躯体化"`
	Obsession     float32 `json:"obsession" gorm:"type:decimal(3,1);not null;comment:强迫症状"`
	Interpersonal float32 `json:"interpersonal" gorm:"type:decimal(3,1);not null;comment:人际关系敏感"`
	Depression    float32 `json:"depression" gorm:"type:decimal(3,1);not null;comment:抑郁"`
	Anxiety       float32 `json:"anxiety" gorm:"type:decimal(3,1);not null;comment:焦虑"`
	Hostility     float32 `json:"hostility" gorm:"type:decimal(3,1);not null;comment:敌对"`
	Phobia        float32 `json:"phobia" gorm:"type:decimal(3,1);not null;comment:恐怖"`
	Paranoia      float32 `json:"paranoia" gorm:"type:decimal(3,1);not null;comment:偏执"`
	Psychoticism  float32 `json:"psychoticism" gorm:"type:decimal(3,1);not null;comment:精神病性"`
	Other         float32 `json:"other" gorm:"type:decimal(3,1);not null;comment:其他"`

	TotalScore    float64 `json:"total_score" gorm:"type:int;default:0;comment:总分（估算值）"`
	PositiveItems float64 `json:"positive_items" gorm:"type:int;default:0;comment:阳性项目数（因子>=2的数量）"`

	CreatedAt time.Time      `json:"-" gorm:"type:timestamp;autoCreateTime;comment:记录创建时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段，代表删除时间
}

// TableName 指定表名为 scl，避免 gorm 自动复数化
func (SCL) TableName() string {
	return "scl"
}
