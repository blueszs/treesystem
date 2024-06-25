package main

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"treesystem/common/image"
	"treesystem/dal"
	"treesystem/models"

	"github.com/gin-gonic/gin"
)

func apiHandler(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "upload":
		upload(c)
	case "uploads":
		uploads(c)
	case "treelist":
		treelist(c)
	case "searchtree":
		treelist(c)
	case "treeinfo":
		treeinfo(c)
	default:
		var result models.Result[any]
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
	}
}

func upload(c *gin.Context) {
	//1、获取上传的文件
	file, err := c.FormFile("files")
	if err != nil {
		var result models.Result[[]string]
		result.Code = 0
		result.Message = "文件类型不合法"
		c.JSON(http.StatusOK, result)
		return
	}
	//2、获取后缀名 判断类型是否正确 .jpg .png .gif .jpeg
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		var result models.Result[[]string]
		result.Code = 0
		result.Message = "文件类型不合法"
		c.JSON(http.StatusOK, result)
		return
	}
	//3、创建图片保存目录,linux下需要设置权限（0666可读可写） /upload/20200623
	day := time.Now().Format("20060102")
	dir := "../view/wwwroot/upload/" + day
	if err := os.MkdirAll(dir, 0666); err != nil {
		//c.String(http.StatusBadRequest, "MkdirAll失败")
		var result models.Result[[]string]
		result.Code = 0
		result.Message = "文件存储目录创建失败"
		c.JSON(http.StatusOK, result)
		return
	}
	//4、生成文件名称 144325235235.png
	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
	//5、上传文件 /upload/20200623/144325235235.png
	var urls []string
	saveDir := path.Join(dir, fileUnixName+extName)
	err = c.SaveUploadedFile(file, saveDir)
	if err != nil {
		var result models.Result[[]string]
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
	} else {
		urls = append(urls, strings.ReplaceAll(saveDir, "../view/wwwroot", ""))
		thumbnailPath := path.Join(dir, fileUnixName+"_thumbnail"+extName) //随机生成新的文件名
		_, err = image.ImageCompress(saveDir, thumbnailPath, 150, 0, 100)
		if err != nil {
			var result models.Result[[]string]
			result.Code = 0
			result.Message = err.Error()
			c.JSON(http.StatusOK, result)
		} else {
			urls = append(urls, strings.ReplaceAll(thumbnailPath, "../view/wwwroot", ""))
			var result models.Result[[]string]
			result.Code = 1
			result.Message = "成功！"
			result.Data = urls
			c.JSON(http.StatusOK, result)
		}
	}
}

func uploads(c *gin.Context) {
	var urls []string
	//1、获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		var result models.Result[[]string]
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	files := form.File["files"]
	//2、获取后缀名 判断类型是否正确 .jpg .png .gif .jpeg
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	for _, file := range files {
		extName := path.Ext(file.Filename)
		if _, ok := allowExtMap[extName]; !ok {
			var result models.Result[[]string]
			result.Code = 0
			result.Message = "文件类型不合法"
			c.JSON(http.StatusOK, result)
			return
		}
		//3、创建图片保存目录,linux下需要设置权限（0666可读可写） /upload/20200623
		day := time.Now().Format("20060102")
		dir := "../view/wwwroot/upload/" + day
		if err := os.MkdirAll(dir, 0666); err != nil {
			//c.String(http.StatusBadRequest, "MkdirAll失败")
			var result models.Result[[]string]
			result.Code = 0
			result.Message = err.Error()
			c.JSON(http.StatusOK, result)
			return
		}
		//4、生成文件名称 144325235235.png
		fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
		//5、上传文件 /upload/20200623/144325235235.png
		saveDir := path.Join(dir, fileUnixName+extName)
		if err := c.SaveUploadedFile(file, saveDir); err != nil {
			var result models.Result[[]string]
			result.Code = 0
			result.Message = err.Error()
			c.JSON(http.StatusOK, result)
			return
		}
		urls = append(urls, strings.ReplaceAll(saveDir, "../view/wwwroot", ""))
	}
	var result models.Result[[]string]
	result.Code = 1
	result.Data = urls
	c.JSON(http.StatusOK, result)
}

func treelist(c *gin.Context) {
	skey := c.PostForm("skey")
	var result models.Result[[]models.TreeInfo]
	mode, err := dal.GetTreeList(sqldb, skey)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = mode
	result.Message = "成功！"
	c.JSON(http.StatusOK, result)
}

func treeinfo(c *gin.Context) {
	sid := c.PostForm("treeid")
	var result models.Result[models.TreeInfo]
	if len(sid) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	mode, err := dal.GetTreeInfo(sqldb, id)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = mode
	result.Message = "成功！"
	c.JSON(http.StatusOK, result)
}
