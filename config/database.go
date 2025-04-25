package config

// 数据库连接配置
import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 负责初始化数据库连接
func InitDB() {
	// 读取配置文件
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 获取数据库配置
	ip := cfg.Section("mysql").Key("ip").String()
	port := cfg.Section("mysql").Key("port").String()
	user := cfg.Section("mysql").Key("user").String()
	password := cfg.Section("mysql").Key("password").String()
	database := cfg.Section("mysql").Key("database").String()

	// 生成 DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, ip, port, database)

	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 记录 SQL 信息
			Colorful:      true,        // 彩色日志输出
		},
	)

	// 连接数据库
	var errDB error
	DB, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if errDB != nil {
		log.Fatalf("数据库连接失败: %v", errDB)
	} else {
		fmt.Println("数据库 连接成功!")
	}
}
