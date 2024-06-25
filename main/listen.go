package main

import (
	"net/http"
	"strings"
	"time"
	"treesystem/dal"
	"treesystem/models"
	"treesystem/sql"

	"github.com/gin-gonic/gin"
)

func listenHandler(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "classinfo":
		classinfoAction(c, action)
	case "seriesinfo":
		seriesinfoAction(c, action)
	case "listeninfo":
		listeninfoAction(c, action)
	default:
		c.HTML(http.StatusOK, "404.html", nil)
	}
}

func classinfoAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetListenInfoList(sqldb)
		return model, "infos", err
	})
}
func seriesinfoAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetListenSeriesList(sqldb)
		return model, "infos", err
	})
}
func listeninfoAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetListenInfoList(sqldb)
		return model, "infos", err
	})
}
func postListen(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "addlistenclass":
		addlistenclass(c)
	case "updatelistenclass":
		updatelistenclass(c)
	case "dellistenclass":
		dellistenclass(c)
	case "addlistenseries":
		addlistenseries(c)
	case "updatelistenseries":
		updatelistenseries(c)
	case "dellistenseries":
		dellistenseries(c)
	case "addlisteninfo":
		addlisteninfo(c)
	case "updatelisteninfo":
		updatelisteninfo(c)
	case "dellisteninfo":
		dellisteninfo(c)
	default:
		var result models.Result[any]
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
	}
}

func addlistenclass(c *gin.Context) {
	var info models.ListenClass
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Class_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	//info.Insert_Time = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddListenClass(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func updatelistenclass(c *gin.Context) {
	var info models.ListenClass
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Class_Id == 0 || len(info.Class_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateListenClass(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func dellistenclass(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelListenClass(sdb, id)
	})
}
func addlistenseries(c *gin.Context) {
	var info models.ListenSeries
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Class_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	info.Add_Time = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddListenSeries(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func updatelistenseries(c *gin.Context) {
	var info models.ListenSeries
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Series_Id == 0 || len(info.Series_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateListenSeries(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func dellistenseries(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelListenSeries(sdb, id)
	})
}
func addlisteninfo(c *gin.Context) {
	var info models.ListenInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Listen_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	info.Add_Time = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddListenInfo(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func updatelisteninfo(c *gin.Context) {
	var info models.ListenInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Listen_Id == 0 || len(info.Listen_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateListenInfo(sqldb, info)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功"
	c.JSON(http.StatusOK, result)
}
func dellisteninfo(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelListenInfo(sdb, id)
	})
}
