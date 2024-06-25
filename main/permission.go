package main

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
	"treesystem/common/encrypt"
	"treesystem/dal"
	"treesystem/models"

	"github.com/gin-gonic/gin"
)

func postPermission(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	switch action {
	case "addconfiginfo":
		addconfiginfo(c)
	case "updateconfiginfo":
		updateconfiginfo(c)
	case "delconfiginfo":
		delconfiginfo(c)
	case "addmoduleinfo":
		addmoduleinfo(c)
	case "updatemoduleinfo":
		updatemoduleinfo(c)
	case "delmoduleinfo":
		delmoduleinfo(c)
	case "addroleinfo":
		addroleinfo(c)
	case "updateroleinfo":
		updateroleinfo(c)
	case "delroleinfo":
		delroleinfo(c)
	case "getrolemodules":
		getrolemodules(c)
	case "updaterolemodule":
		updaterolemodule(c)
	case "adduserinfo":
		adduserinfo(c)
	case "updateuserinfo":
		updateuserinfo(c)
	case "deluserinfo":
		deluserinfo(c)
	default:
		var result models.Result[any]
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
	}
}

func addconfiginfo(c *gin.Context) {
	var info models.SystemConfig
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Config_Key) == 0 {
		result.Code = 0
		result.Message = "配置参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.AddSystemConfig(sqldb, info)
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

func updateconfiginfo(c *gin.Context) {
	var info models.SystemConfig
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Old_Config_Key) == 0 || len(info.Config_Key) == 0 {
		result.Code = 0
		result.Message = "配置参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateSystemConfig(sqldb, info)
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

func delconfiginfo(c *gin.Context) {
	configKey := c.PostForm("configKey")
	var result models.Result[bool]
	if len(configKey) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.DelSystemConfig(sqldb, configKey)
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

func addmoduleinfo(c *gin.Context) {
	var info models.ModuleInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Module_Name) == 0 || len(info.Module_Url) == 0 || len(info.Module_Icon) == 0 || len(info.Module_Info) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	info.Module_Addtime = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddModuleInfo(sqldb, info)
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
func updatemoduleinfo(c *gin.Context) {
	var info models.ModuleInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Module_Id == 0 || len(info.Module_Name) == 0 || len(info.Module_Url) == 0 || len(info.Module_Icon) == 0 || len(info.Module_Info) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateModuleInfo(sqldb, info)
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
func delmoduleinfo(c *gin.Context) {
	moduleId := c.PostForm("moduleId")
	var result models.Result[bool]
	if len(moduleId) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	id, err := strconv.Atoi(moduleId)
	if err != nil {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.DelModuleInfo(sqldb, id)
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
func getrolemodules(c *gin.Context) {
	roleId := c.PostForm("roleId")
	var result models.Result[[]models.ModuleInfo]
	if len(roleId) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	id, err := strconv.Atoi(roleId)
	if err != nil {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	modes, err := dal.GetRoleModules(sqldb, id)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	result.Code = 1
	result.Data = modes
	result.Message = "成功！"
	c.JSON(http.StatusOK, result)
}
func updaterolemodule(c *gin.Context) {
	var info models.RoleModule
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.RoleId == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateRoleModule(sqldb, info)
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

func addroleinfo(c *gin.Context) {
	var info models.RoleInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.Role_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	info.Insert_Time = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddRoleInfo(sqldb, info)
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
func updateroleinfo(c *gin.Context) {
	var info models.RoleInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Role_Id == 0 || len(info.Role_Name) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.UpdateRoleInfo(sqldb, info)
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
func delroleinfo(c *gin.Context) {
	roleId := c.PostForm("roleId")
	var result models.Result[bool]
	if len(roleId) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	id, err := strconv.Atoi(roleId)
	if err != nil {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.DelRoleInfo(sqldb, id)
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
func adduserinfo(c *gin.Context) {
	var info models.UserInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if len(info.User_Id) == 0 || len(info.User_Password) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	pwddata := []byte(info.User_Password)
	keydata := []byte(Conf.System.PwdKey)
	pwdenc, err := encrypt.AesECBEncrypt(pwddata, keydata)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	info.User_Password = hex.EncodeToString(pwdenc)
	info.Insert_Time = time.Now().Format("2006-01-02 15:04:05")
	bl, err := dal.AddUserInfo(sqldb, info)
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
func updateuserinfo(c *gin.Context) {
	var info models.UserInfo
	c.BindJSON(&info)
	var result models.Result[bool]
	if info.Id == 0 || len(info.User_Id) == 0 || len(info.User_Password) == 0 {
		result.Code = 0
		result.Message = "参数值错误"
		c.JSON(http.StatusOK, result)
		return
	}
	pwddata := []byte(info.User_Password)
	keydata := []byte(Conf.System.PwdKey)
	pwdenc, err := encrypt.AesECBEncrypt(pwddata, keydata)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	info.User_Password = hex.EncodeToString(pwdenc)
	bl, err := dal.UpdateUserInfo(sqldb, info)
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
func deluserinfo(c *gin.Context) {
	userId := c.PostForm("userId")
	var result models.Result[bool]
	if len(userId) == 0 {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		result.Code = 0
		result.Message = "参数错误"
		c.JSON(http.StatusOK, result)
		return
	}
	bl, err := dal.DelUserInfo(sqldb, id)
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
