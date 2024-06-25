package common

import (
	"math/rand"
	"time"
)

const numberAndAlphabet = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

// GetRandomStr 获取随机字符串
// @randStrLen 随机字符串长度
func GetRandomStr(randStrLen int) string {
	var result string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rest = []rune(numberAndAlphabet)
	for i := 0; i < randStrLen; i++ {
		j := r.Intn(len(numberAndAlphabet))
		result += string(rest[j])
	}
	return result
}
