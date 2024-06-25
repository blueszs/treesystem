package main

import (
	"net/http"
	"strconv"
	"strings"
	"treesystem/dal"
	"treesystem/models"
	"treesystem/sql"

	"github.com/gin-gonic/gin"
)

func treeHandler(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "treeplant":
		treeplantAction(c, action)
	case "treeclass":
		treeclassAction(c, action)
	case "treeinfo":
		treeinfoAction(c, action)
	case "treedetail":
		treedetailAction(c, action)
	default:
		c.HTML(http.StatusOK, "404.html", nil)
	}
}

func treeplantAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetTreePlantList(sqldb)
		return model, "infos", err
	})
}

func treeclassAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetTreeClassList(sqldb)
		return model, "infos", err
	})
}
func treeinfoAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetTreeInfoList(sqldb)
		return model, "infos", err
	}, func() (any, string, error) {
		model, err := dal.GetTreeClassList(sqldb)
		return model, "classes", err
	})
}
func treedetailAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetTreeDetailList(sqldb)
		return model, "infos", err
	}, func() (any, string, error) {
		model, err := dal.GetTreeInfoList(sqldb)
		return model, "classes", err
	}, func() (any, string, error) {
		model, err := dal.GetTreePlantList(sqldb)
		return model, "plants", err
	})
}

func postTree(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "addtreeclass":
		addtreeclass(c)
	case "updatetreeclass":
		updatetreeclass(c)
	case "deltreeclass":
		deltreeclass(c)
	case "addtreeplant":
		addtreeplant(c)
	case "updatetreeplant":
		updatetreeplant(c)
	case "deltreeplant":
		deltreeplant(c)
	case "addtreeinfo":
		addtreeinfo(c)
	case "updatetreeinfo":
		updatetreeinfo(c)
	case "deltreeinfo":
		deltreeinfo(c)
	case "gettreedetail":
		gettreedetail(c)
	case "addtreedetail":
		addtreedetail(c)
	case "updatetreedetail":
		updatetreedetail(c)
	case "deltreedetail":
		deltreedetail(c)
	default:
		var result models.Result[any]
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
	}
}
func addtreeclass(c *gin.Context) {
	var info models.TreeClass
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Class_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.AddTreeClass(sqldb, info)
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
func updatetreeclass(c *gin.Context) {
	var info models.TreeClass
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Class_Id == 0 || len(info.Class_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateTreeClass(sqldb, info)
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
func deltreeclass(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelTreeClass(sdb, id)
	})
}
func addtreeplant(c *gin.Context) {
	var info models.TreePlant
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Plant_Name) == 0 || len(info.Plant_Address) == 0 || len(info.Plant_Contacts) == 0 || len(info.Plant_Tel) == 0 || len(info.Plant_Info) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.AddTreePlant(sqldb, info)
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
func updatetreeplant(c *gin.Context) {
	var info models.TreePlant
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Plant_Id == 0 || len(info.Plant_Name) == 0 || len(info.Plant_Address) == 0 || len(info.Plant_Contacts) == 0 || len(info.Plant_Tel) == 0 || len(info.Plant_Info) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateTreePlant(sqldb, info)
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
func deltreeplant(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelTreePlant(sdb, id)
	})
}

func addtreeinfo(c *gin.Context) {
	var info models.TreeInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Class_Id == 0 || len(info.Tree_Name) == 0 || len(info.Tree_Info) == 0 || len(info.Tree_Main_Photo) == 0 || len(info.Tree_Main_Photo_Thumbnail) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.AddTreeInfo(sqldb, info)
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
func updatetreeinfo(c *gin.Context) {
	var info models.TreeInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Tree_Id == 0 || info.Class_Id == 0 || len(info.Tree_Name) == 0 || len(info.Tree_Info) == 0 || len(info.Tree_Main_Photo) == 0 || len(info.Tree_Main_Photo_Thumbnail) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateTreeInfo(sqldb, info)
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
func deltreeinfo(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelTreeInfo(sdb, id)
	})
}
func gettreedetail(c *gin.Context) {
	sid := c.PostForm("id")
	var result models.Result[models.TreeDetail]
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
	mode, err := dal.GetTreeDetail(sqldb, id)
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

func addtreedetail(c *gin.Context) {
	var info models.TreeDetail
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Tree_Id == 0 || info.Plant_Id == 0 || info.Active_Price == 0 || info.Real_Price == 0 || info.Prime_Price == 0 || len(info.Tree_Unity) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.AddTreeDetail(sqldb, info)
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
func updatetreedetail(c *gin.Context) {
	var info models.TreeDetail
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Detail_Id == 0 || info.Tree_Id == 0 || info.Plant_Id == 0 || info.Active_Price == 0 || info.Real_Price == 0 || info.Prime_Price == 0 || len(info.Tree_Unity) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateTreeDetail(sqldb, info)
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
func deltreedetail(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelTreeDetail(sdb, id)
	})
}

func delfunc(c *gin.Context, f func(*sql.MySqlDb, int) (bool, error)) {
	sid := c.PostForm("id")
	var result models.Result[bool]
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
	bl, err := f(sqldb, id)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = bl
	result.Message = "成功！"
	c.JSON(http.StatusOK, result)
}
