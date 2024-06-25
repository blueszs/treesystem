package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// DeferFileCols 文件流关闭
func DeferFileCols(f *os.File) {
	if f != nil {
		err := f.Close()
		if err != nil {
			fmt.Printf("file close failed:%v\n", err)
		} else {
			fmt.Println("file close success!")
		}
	} else {
		fmt.Println("file object is null or undefined!")
	}
}

// Save 将内存的数据文件存储到系统的文件夹中
// @stream 内存字节流
// @sPath 存储路径
// @coverage 是否覆盖现有文件
func Save(stream []byte, sPath string, coverage bool) (bool, error) {
	dirs := strings.Split(sPath, "\\")
	//获取文件存储的路径
	vPath := strings.Join(dirs[0:len(dirs)-1], "\\")
	//判断路径是否存在
	bl, err := PathExists("", vPath)
	if err != nil {
		return false, err
	}
	//路径不存在时先创建路径
	if !bl {
		//创建文件存储路径
		bl, err = CreatePathDir(vPath)
		if err != nil {
			return bl, err
		}
	}
	//如果不覆盖文件需要判断文件是否存在
	if !coverage {
		//判断文件是否存在
		bl, err := PathExists("", sPath)
		if err != nil {
			return false, err
		}
		if bl {
			return false, errors.New("file is exists")
		}
	}
	//将内存中的字节流写入文件
	err = os.WriteFile(sPath, stream, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

// PathExists 判断路径是否存在
// @startPath 起始路径
// @sPath 路径地址
func PathExists(startPath, sPath string) (bool, error) {
	if startPath == "" {
		dirs := strings.Split(sPath, "\\")
		for i := 0; i < len(dirs); i++ {
			if i == 0 {
				// 判断根目录是否存在
				flag, err := PathExists(dirs[0], "")
				if err == nil && !flag {
					return false, errors.New("root path is not exists")
				}
			} else {
				// 判断根目录是否存在
				flag, err := PathExists(strings.Join(dirs[:i], "\\"), dirs[i])
				if err == nil && !flag {
					return false, nil
				}
			}
		}
	} else {
		_, err := os.Stat(startPath + "\\" + sPath)
		if err == nil {
			return true, nil
		} else {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
	}
	return true, nil
}

// CreatePathDir 创建所有路径
// @sPath 路径地址
func CreatePathDir(sPath string) (bool, error) {
	dirs := strings.Split(sPath, "\\")
	if len(dirs) > 1 {
		// 判断根目录是否存在
		bl, err := PathExists(dirs[0], "")
		if err != nil {
			return false, err
		} else {
			if bl {
				//创建不存在的所有路径
				err = os.MkdirAll(sPath, 0766)
				if err != nil {
					return false, err
				} else {
					return true, nil
				}
			} else {
				return false, errors.New("root path is not exists")
			}
		}
	}
	return false, errors.New("'sPath' params is not null")
}

// ReadFile 根据文件路径读取文件内容返回内容字符串
// @sPath 文件路径
func ReadFile(sPath string) (string, error) {
	//判断路径是否存在
	_, err := PathExists("", sPath)
	if err != nil {
		return "", err
	}
	buff, err := os.ReadFile(sPath)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
