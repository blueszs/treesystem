package dal

import (
	"strconv"
	"treesystem/models"
	"treesystem/sql"
)

// GetUserInfo 根据关键字搜索，获取树种列表
// @mysql 数据库实例结构体
// @username 用户名
// @password 密码
func GetUserInfo(mysql *sql.MySqlDb, username, password string) (models.UserInfo, error) {
	sqlstr := "select acc.*,GROUP_CONCAT(ro.role_id) as role_id,GROUP_CONCAT(ro.role_name) as role_name FROM user_info acc,role_info ro ,user_role accRole WHERE acc.user_id=? and acc.user_password=? AND (acc.id = accRole.user_id) AND accRole.role_id = ro.role_id GROUP BY accRole.user_id ORDER BY acc.id "
	query, err := mysql.QuerySql(sqlstr, username, password)
	if err != nil {
		return *new(models.UserInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.UserInfo, error) {
		var mo models.UserInfo
		mo.Id, err = strconv.Atoi(kr["id"])
		if err != nil {
			return *new(models.UserInfo), err
		}
		mo.User_Id = kr["user_id"]
		mo.User_Name = kr["user_name"]
		mo.User_Password = kr["user_password"]
		mo.Insert_Time = kr["insert_time"]
		mo.RoleId = kr["role_id"]
		mo.RoleName = kr["role_name"]
		return mo, nil
	})
	if err != nil {
		return *new(models.UserInfo), err
	}
	return data[0], nil
}

// GetUserNavbars 获取用户主菜单权限
// @userid 用户编号
func GetUserNavbars(mysql *sql.MySqlDb, userid int) ([]models.UserNavbar, error) {
	sqlstr := "select b.role_id,b.role_name,c.module_id,d.module_name,d.module_url,d.parent_module_id,GROUP_CONCAT(c.module_operation)as module_operation,d.module_icon from user_role a,role_info b ,role_module c ,module_info d where a.role_id= b.role_id and c.module_id = d.module_id and a.role_id = c.role_id  and a.user_id=? GROUP BY c.module_id"
	query, err := mysql.QuerySql(sqlstr, userid)
	if err != nil {
		return *new([]models.UserNavbar), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.UserNavbar, error) {
		var mo models.UserNavbar
		mo.RoleId, err = strconv.Atoi(kr["role_id"])
		if err != nil {
			return *new(models.UserNavbar), err
		}
		mo.RoleName = kr["role_name"]
		mo.ModuleId, err = strconv.Atoi(kr["module_id"])
		if err != nil {
			return *new(models.UserNavbar), err
		}
		mo.ModuleName = kr["module_name"]
		mo.ModuleUrl = kr["module_url"]
		mo.ParentModuleId, err = strconv.Atoi(kr["parent_module_id"])
		if err != nil {
			return *new(models.UserNavbar), err
		}
		mo.ModuleOperation = kr["module_operation"]
		mo.ModuleIcon = kr["module_icon"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.UserNavbar), err
	}
	return data, nil
}
