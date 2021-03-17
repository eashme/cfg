package coord_cfg

import (
	"encoding/json"
	"log"
	"strings"
)

// 配置文件 code 的标准命名,存DB和程序运行中都存储的该命名
func StandCode(code string) string {
	return strings.TrimSpace( // 去除空格
		strings.ToUpper( // 转大写
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(code,
						"/", "_"), // 转 / 成 _
					".", "_"), // .转成 _
				"-", "_"))) // 把 - 转成 _
}

func Json(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}


func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error ", err.Error())
	}
}
