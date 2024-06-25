package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"treesystem/common"
	"treesystem/common/encrypt"
	"treesystem/dal"
	"treesystem/models"

	"github.com/gin-gonic/gin"
)

// Cors 定义全局的中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func autokenHandler(f func(*gin.Context, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		action := c.Param("action")
		if strings.ToLower(action) == "login" {
			c.Next()
		} else {
			auth, err := c.Request.Cookie("authtoken")
			if err != nil {
				f(c, err)
				return
			}
			if CacheDb.Flag {
				token := CacheDb.Get("bk:user:" + common.GetStringMD5(auth.Value))
				if len(token) == 0 {
					return
				}
				userInfo, err := common.JsonDeserialize[models.UserInfo](token)
				if err != nil {
					f(c, err)
					return
				}
				c.Set("userInfo", userInfo)
				c.Next()

			} else {
				tkeydata := []byte(Conf.System.TokenKey)
				v, err := hex.DecodeString(auth.Value)
				if err != nil {
					f(c, err)
					return
				}
				pwdenc, err := encrypt.AesECBDecrypt(v, tkeydata)
				if err != nil {
					f(c, err)
					return
				}
				sjson := string(pwdenc)
				userInfo, err := common.JsonDeserialize[models.UserInfo](sjson)
				if err != nil {
					f(c, err)
					return
				}
				c.Set("userInfo", userInfo)
				c.Next()
			}
		}
	}
}

func userBarsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := c.Param("action")
		if strings.ToLower(action) == "login" {
			c.Next()
		} else {
			userInfo, bl := c.Get("userInfo")
			if !bl {
				c.HTML(500, "500.html", gin.H{})
				c.Abort()
				return
			}
			m := userInfo.(models.UserInfo)
			//判断是否配置缓存，是否重缓存中读取
			if CacheDb.Flag {
				ckey := "bk:user:navbar:" + strconv.Itoa(m.Id)
				redstr := CacheDb.Get(ckey)
				if len(redstr) == 0 {
					userbars, err := dal.GetUserNavbars(sqldb, m.Id)
					if err != nil {
						c.HTML(500, "500.html", gin.H{})
						c.Abort()
						return
					}
					c.Set("Userbars", userbars)
					v, err := json.Marshal(&userbars)
					if err != nil {
						c.HTML(500, "500.html", gin.H{})
						c.Abort()
						return
					}
					CacheDb.Set(ckey, v, 0)
				} else {
					userbars, err := common.JsonDeserialize[[]models.UserNavbar](redstr)
					if err != nil {
						c.HTML(500, "500.html", gin.H{})
						c.Abort()
						return
					}
					c.Set("Userbars", userbars)
				}
			} else {
				userbars, err := dal.GetUserNavbars(sqldb, m.Id)
				if err != nil {
					c.HTML(500, "500.html", gin.H{})
					c.Abort()
					return
				}
				c.Set("Userbars", userbars)
			}
			c.Next()
		}
	}
}

func operatHandler(f func(*gin.Context, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, bl := c.Get("userInfo")
		if !bl {
			f(c, errors.New("获取用户信息失败"))
			return
		}
		m := userInfo.(models.UserInfo)
		userbars, err := dal.GetUserNavbars(sqldb, m.Id)
		if err != nil {
			f(c, err)
			return
		}
		var checked bool
		operatAuth := make([]string, 0)
		for _, v := range userbars {
			if v.ModuleUrl == c.Request.RequestURI {
				if !checked {
					checked = true
				}
				if len(v.ModuleOperation) > 0 {
					for _, o := range v.ModuleOperation {
						var operflag bool
						for _, v2 := range operatAuth {
							if v2 == string(o) {
								operflag = true
								break
							}
						}
						if !operflag {
							operatAuth = append(operatAuth, string(o))
						}
					}
				}
			}
		}
		if checked {
			c.Set("Operat", operatAuth)
			c.Next()
		} else {
			f(c, errors.New("用户无操作权限"))
		}
	}
}

func OperatCheckHandler(f func(*gin.Context, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, bl := c.Get("userInfo")
		if !bl {
			f(c, errors.New("获取用户信息失败"))
			return
		}
		m := userInfo.(models.UserInfo)
		userbars, err := dal.GetUserNavbars(sqldb, m.Id)
		if err != nil {
			f(c, err)
			return
		}
		var checked bool
		operatAuth := make([]string, 0)
		for _, v := range userbars {
			if v.ModuleUrl == c.Request.RequestURI {
				if !checked {
					checked = true
				}
				if len(v.ModuleOperation) > 0 {
					for _, o := range v.ModuleOperation {
						var operflag bool
						for _, v2 := range operatAuth {
							if v2 == string(o) {
								operflag = true
								break
							}
						}
						if !operflag {
							operatAuth = append(operatAuth, string(o))
						}
					}
				}
			}
		}
		if checked {
			c.Set("Operat", operatAuth)
			c.Next()
		} else {
			f(c, errors.New("用户无操作权限"))
		}
	}
}
