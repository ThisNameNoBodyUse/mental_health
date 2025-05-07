package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mental/config"
	"mental/models"
	"os/exec"
	"strings"
)

type Swagger struct {
	Paths map[string]map[string]struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
	} `json:"paths"`
}

// InsertSwaggerAPIs 解析 swagger.json 并插入数据库
func InsertSwaggerAPIs(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取 swagger 文件失败: %v", err)
	}

	var swagger Swagger
	if err := json.Unmarshal(data, &swagger); err != nil {
		return fmt.Errorf("解析 swagger JSON 失败: %v", err)
	}

	var apis []models.API
	for path, methods := range swagger.Paths {
		for method, info := range methods {
			api := models.API{
				Path:        path,
				Method:      strings.ToUpper(method),
				Description: info.Summary + " - " + info.Description,
			}
			apis = append(apis, api)
		}
	}

	if len(apis) > 0 {
		result := config.DB.Create(&apis)
		if result.Error != nil {
			return fmt.Errorf("插入数据库失败: %v", result.Error)
		}
		fmt.Printf("成功插入 %d 条 API 接口\n", len(apis))
	}
	return nil
}

// RunSwagInit 自动运行 swag init 命令
func RunSwagInit() error {
	cmd := exec.Command("swag", "init") // swag 命令必须在环境变量中
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("swag init 执行失败: %v\n输出: %s", err, string(output))
	}
	fmt.Println("swag init 执行成功")
	return nil
}
