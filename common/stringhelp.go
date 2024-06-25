package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
)

type SStr string

// ConvertUnicodeToChar 转换Unicode字符串到中文
// @in 输入Unicode字符串
// eg: `\u82f9\u679c`
func ConvertUnicodeToChar(in string) string {
	reg := regexp.MustCompile(`\\u([0-9]|[a-z]|[A-Z]){4}`)
	out := reg.ReplaceAllStringFunc(in,
		func(b string) string {
			v := strings.ReplaceAll(b, "\\u", "")
			temp, err := strconv.ParseInt(v, 16, 32)
			if err != nil {
				fmt.Println("convert error ", err.Error())
				panic(err)
			}
			i := fmt.Sprintf("%c", temp)
			return i
		})
	return out
}

// ConvertCharToUnicode 获取字符串的Unicode编码
// @in 输入中文字符串
// eg 输入‘苹果’输出`\u82f9\u679c`
func ConvertCharToUnicode(in string) string {
	code := strconv.QuoteToASCII(in)
	return code[1 : len(code)-1]
}

// FullToUnicode 将字符串全部转化为Unicode字符编码（包含字母、数字、符号等）
// @in 待转码文字符串
func FullToUnicode(in string) string {
	DD := []rune(in) //需要分割的字符串内容，将它转为字符，然后取长度。
	finalStr := ""
	for i := 0; i < len(DD); i++ {
		if unicode.Is(unicode.Scripts["Han"], DD[i]) {
			textQuoted := strconv.QuoteToASCII(string(DD[i]))
			finalStr += textQuoted[1 : len(textQuoted)-1]
		} else {
			h := fmt.Sprintf("%x", DD[i])
			switch expr := len(h); expr {
			case 1:
				finalStr += "\\u" + "000" + h
			case 2:
				finalStr += "\\u" + "00" + h
			case 3:
				finalStr += "\\u" + "0" + h
			default:
				finalStr += "\\u" + h
			}
		}
	}
	return finalStr
}

// ReplaceByRegex 根据正则表达式替换字符串
// @register 正则表达式
// @oldster 需要替换的旧字符串
// @rest 替换的新字符串
func ReplaceByRegex(register, oldster, rest string) string {
	reg := regexp.MustCompile(register)
	var out = reg.ReplaceAllStringFunc(oldster,
		func(b string) string {
			v := strings.ReplaceAll(b, b, rest)
			return v
		})
	return out
}

// GetRegexValue 根据正则表达式获取字符串
// @register 正则表达式
// @instr 需要匹配的字符串
func GetRegexValue(register, instr string) string {
	reg, _ := regexp2.Compile(register, 0)
	//reg := regexp.MustCompile(register)
	out, _ := reg.FindStringMatch(instr)
	return out.String()
}

// AppendByte 字符切片添加新的字节
// @slice 原始字节切片
// @data 新增的字节
func AppendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // 如果需要，重新分配
		// 为未来的增长分配双倍的需求。
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

// Utf8Decode 字符解密
// @data 待解密字节切片
func Utf8Decode(data []byte) string {
	fmt.Printf("%c", data)
	fmt.Println("")
	fmt.Printf("%s", string(data))
	fmt.Println("")
	//var result string
	tData := make([]rune, 0, len(data)+1)
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		tData = append(tData, r)
		data = data[size:]
	}
	return string(tData)
}

// Utf8Encode 字符解密
// @str 待加密字节切片
func Utf8Encode(str string) string {
	rData := []rune(str)
	//var result string
	tData := make([]byte, 0, len(rData)*3)
	for _, r := range rData {
		tmpD := make([]byte, 3)
		utf8.EncodeRune(tmpD, r)
		tData = append(tData, tmpD...)
	}
	return string(tData)
}

func (s SStr) SubString(start int, end int) string {
	var r = []rune(s)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return string(r)
	}
	return string(r[start:end])
}

// GetStringMD5 获取字符串MD5码
func GetStringMD5(str string) string {

	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
