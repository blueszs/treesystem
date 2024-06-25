package common

import (
	"encoding/json"
)

// ReadConfig 读取配置文件内容
// @spath 配置文件路径
func ReadConfig(spath string) (string, error) {
	str, err := ReadFile(spath)
	if err != nil {
		return "", err
	}
	return str, nil
}

// JsonDeserialize 将json字符串系列化为实例
// @sjson json字符串
func JsonDeserialize[T any](sjson string) (T, error) {
	var mode T
	err := json.Unmarshal([]byte(sjson), &mode)
	if err != nil {
		return *new(T), err
	}
	return mode, nil
}
