package main

import (
	"net/http"
	"strings"
	"treesystem/dal"
	"treesystem/models"
	"treesystem/sql"

	"github.com/gin-gonic/gin"
)

func memberHandler(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "memberlist":
		memberlistAction(c, action)
	default:
		c.HTML(http.StatusOK, "404.html", nil)
	}
}

func memberlistAction(c *gin.Context, action string) {
	permissionAction(c, action, func() (any, string, error) {
		model, err := dal.GetMemberInfoList(sqldb)
		return model, "infos", err
	})
}

func postMember(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "delmemberinfo":
		delmemberinfo(c)
	default:
		var result models.Result[any]
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
	}
}

func delmemberinfo(c *gin.Context) {
	delfunc(c, func(sdb *sql.MySqlDb, id int) (bool, error) {
		return dal.DelMemberInfo(sdb, id)
	})
}
