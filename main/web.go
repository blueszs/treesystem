package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"treesystem/common"
	"treesystem/common/encrypt"
	"treesystem/dal"
	"treesystem/models"
	"treesystem/sql"

	"github.com/gin-gonic/gin"
)

var sqldb *sql.MySqlDb

func RunWebService(sql *sql.MySqlDb) {

	sqldb = sql
	//运行模式
	gin.SetMode(gin.ReleaseMode)
	//Default返回一个默认的路由引擎
	router := gin.Default()
	// 静态资源加载，本例为css,js以及资源图片
	//router.Static 指定某个目录为静态资源目录，可直接访问这个目录下的资源，url 要具体到资源名称。
	//router.StaticFS 比前面一个多了个功能，当目录下不存 index.html 文件时，会列出该目录下的所有文件。
	//router.StaticFile 指定某个具体的文件作为静态资源访问。
	router.Static("/script", "../view/wwwroot/script")
	router.Static("/style", "../view/wwwroot/css")
	router.Static("/fonts", "../view/wwwroot/font")
	router.Static("/webfonts", "../view/wwwroot/webfont")
	router.Static("/lib", "../view/wwwroot/lib")
	router.Static("/upload", "../view/wwwroot/upload")
	router.Static("/resources", "../view/wwwroot/resources/")
	// 导入所有模板，多级目录结构需要这样写
	router.LoadHTMLGlob("../view/tpl/*")

	// website分组
	v := router.Group("/")
	{
		v.GET("/:action", autokenHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), userBarsHandler(), actionHandler)
		v.POST("/authorization", authorizationHandler)
		/*
			v.GET("/user/:name/*action", func(c *gin.Context) {
				name := c.Param("name")
				action := c.Param("action")
				message := name + " is " + action
				c.String(http.StatusOK, message)
			})*/
		v.POST("/form_post", func(c *gin.Context) {
			message := c.PostForm("message")
			nick := c.DefaultPostForm("nick", "anonymous")
			c.JSON(http.StatusOK, gin.H{
				"status": gin.H{
					"status_code": http.StatusOK,
					"status":      "ok",
				},
				"message": message,
				"nick":    nick,
			})
		})
		v.POST("/upload", func(c *gin.Context) {
			name := c.PostForm("name")
			fmt.Println(name)
			file, header, err := c.Request.FormFile("upload")
			if err != nil {
				c.String(http.StatusBadRequest, "Bad request")
				return
			}
			filename := header.Filename
			fmt.Println(file, err, filename)
			out, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				log.Fatal(err)
			}
			c.String(http.StatusCreated, "upload successful")
		})
	}
	permission := router.Group("/permission")
	{
		permission.GET("/:action", autokenHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{"error": err.Error()})
			c.Abort()
		}), operatHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{"error": err.Error()})
			c.Abort()
		}), actionHandler)
		permission.POST("/:action", autokenHandler(func(c *gin.Context, err error) {
			var mo models.Result[any]
			mo.Code = 0
			mo.Message = "未经授权的访问"
			c.JSON(http.StatusOK, mo)
			c.Abort()
		}), postPermission)
	}
	tree := router.Group("/tree")
	{
		tree.GET("/:action", autokenHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), operatHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), treeHandler)
		tree.POST("/:action", autokenHandler(func(c *gin.Context, err error) {
			var mo models.Result[any]
			mo.Code = 0
			mo.Message = "未经授权的访问"
			c.JSON(http.StatusOK, mo)
			c.Abort()
		}), postTree)
	}
	member := router.Group("/member")
	{
		member.GET("/:action", autokenHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), operatHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), memberHandler)
		member.POST("/:action", autokenHandler(func(c *gin.Context, err error) {
			var mo models.Result[any]
			mo.Code = 0
			mo.Message = "未经授权的访问"
			c.JSON(http.StatusOK, mo)
			c.Abort()
		}), postMember)
	}
	listen := router.Group("/listen")
	{
		listen.GET("/:action", autokenHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), operatHandler(func(c *gin.Context, err error) {
			c.HTML(200, "500.html", gin.H{})
			c.Abort()
		}), listenHandler)
		listen.POST("/:action", autokenHandler(func(c *gin.Context, err error) {
			var mo models.Result[any]
			mo.Code = 0
			mo.Message = "未经授权的访问"
			c.JSON(http.StatusOK, mo)
			c.Abort()
		}), postListen)
	}
	// API分组(RESTFULL)以及版本控制
	v1 := router.Group("/v1")
	{
		//下面是群组中间的用法
		v1.POST("/api/:action", apiHandler)
	}
	vport := fmt.Sprintf(":%s", strconv.Itoa(Conf.Web.Port))
	router.Run(vport)
	/*srv := &http.Server{
		Addr:    vport,
		Handler: router,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")*/
}

func actionHandler(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "login":
		c.HTML(http.StatusOK, "login.html", nil)
	case "home":
		homeAction(c, action)
	case "configlist":
		configlistAction(c, action)
	case "modulelist":
		modulelistAction(c, action)
	case "rolelist":
		rolelistAction(c, action)
	case "rolemodule":
		rolemoduleAction(c, action)
	case "userlist":
		userlistAction(c, action)
	default:
		c.HTML(http.StatusOK, "404.html", nil)
	}
}

func getPageOperat(operat []string) models.PageOperat {
	var m models.PageOperat
	for _, v := range operat {
		switch v {
		case "I":
			m.IsAdd = 1
		case "D":
			m.IsDel = 1
		case "S":
			m.IsQuery = 1
		case "U":
			m.IsUpdate = 1
		}
	}
	return m
}

func homeAction(c *gin.Context, action string) {
	userBars, _ := c.Get("Userbars")
	userInfo, _ := c.Get("userInfo")
	bars := userBars.([]models.UserNavbar)
	parent := make([]models.UserNavbar, 0)
	for _, v := range bars {
		if v.ParentModuleId == 0 {
			var count int
			for _, sn := range bars {
				if v.ModuleId == sn.ParentModuleId {
					count += 1
					v.SonModule = append(v.SonModule, sn)
				}
			}
			v.SonModuleFlag = count > 0
			parent = append(parent, v)
		}
	}
	m := userInfo.(models.UserInfo)
	c.HTML(http.StatusOK, action+".html", gin.H{
		"bars": parent,
		"user": m,
	})
}

func configlistAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetSystemConfigList(sqldb)
		return model, "infos", err
	})
}

func modulelistAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetModuleInfoList(sqldb)
		return model, "infos", err
	})
}

func rolelistAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetRoleInfoList(sqldb)
		return model, "infos", err
	})
}

func rolemoduleAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetRoleInfoList(sqldb)
		return model, "roles", err
	}, func() (any, string, error) {
		model, err := dal.GetModuleInfoList(sqldb)
		return model, "modules", err
	})
}
func userlistAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetUserInfoList(sqldb)
		return model, "infos", err
	}, func() (any, string, error) {
		model, err := dal.GetRoleInfoList(sqldb)
		return model, "roles", err
	})
}

func permissionAction(c *gin.Context, action string, f ...func() (any, string, error)) {
	userInfo, _ := c.Get("userInfo")
	m := userInfo.(models.UserInfo)
	pageOperat, _ := c.Get("Operat")
	opera := pageOperat.([]string)
	pageOper := getPageOperat(opera)
	model := make(map[string]any, 0)
	for _, v := range f {
		models, str, err := v()
		if err != nil {
			log.Fatal("sql err:", err)
			c.String(http.StatusServiceUnavailable, "Bad request")
			return
		}
		model[str] = models
	}
	c.HTML(http.StatusOK, action+".html", gin.H{
		"operat": pageOper,
		"user":   m,
		"models": model,
	})
}

func authorizationHandler(c *gin.Context) {
	username := c.PostForm("username")
	var result models.Result[models.UserInfo]
	if username == "" {
		result.Code = 0
		result.Message = "用户名不能为空"
		c.JSON(http.StatusOK, result)
		return
	}
	password := c.PostForm("password")
	if password == "" {
		result.Code = 0
		result.Message = "密码不能为空"
		c.JSON(http.StatusOK, result)
		return
	}
	pwddata := []byte(password)
	keydata := []byte(Conf.System.PwdKey)
	pwdenc, err := encrypt.AesECBEncrypt(pwddata, keydata)
	if err != nil {
		log.Fatal("aes encrypt err:", err)
	}
	model, err := dal.GetUserInfo(sqldb, username, hex.EncodeToString(pwdenc))
	if err != nil {
		log.Fatal("sql err:", err)
	}
	model.User_Password = ""
	v, err := json.Marshal(&model)
	if err != nil {
		log.Fatal("convert json err:", err)
	}
	tkeydata := []byte(Conf.System.TokenKey)
	pwdenc, err = encrypt.AesECBEncrypt(v, tkeydata)
	if err != nil {
		log.Fatal("aes encrypt err:", err)
	}
	model.Token = hex.EncodeToString(pwdenc)
	if CacheDb.Flag {
		CacheDb.Set("bk:user:"+common.GetStringMD5(model.Token), v, 3600)
	}
	result.Code = 1
	result.Data = model
	c.JSON(http.StatusOK, result)
}
