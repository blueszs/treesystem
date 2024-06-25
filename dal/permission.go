package dal

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"treesystem/models"
	"treesystem/sql"
)

// GetSystemConfigList 获取所有配置在数据库中的系统参数
// @mysql 数据库实例
func GetSystemConfigList(mysql *sql.MySqlDb) ([]models.SystemConfig, error) {
	sqlstr := "select a.* from system_config a"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.SystemConfig), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.SystemConfig, error) {
		var mo models.SystemConfig
		mo.Config_Key = kr["config_key"]
		mo.Config_Value = kr["config_value"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.SystemConfig), err
	}
	return data, nil
}

// AddSystemConfig 在数据库中新增系统参数配置
// @mysql 数据库实例
// @info 参数配置实例
func AddSystemConfig(mysql *sql.MySqlDb, info models.SystemConfig) (bool, error) {
	sqlstr := "insert into system_config(config_key,config_value)values(?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Config_Key, info.Config_Value)
	return bl, err
}

// UpdateSystemConfig 修改数据库中的系统参数配置
// @mysql 数据库实例
// @info 参数配置实例
func UpdateSystemConfig(mysql *sql.MySqlDb, info models.SystemConfig) (bool, error) {
	sqlstr := "update system_config set config_key=?,config_value=? where config_key=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Config_Key, info.Config_Value, info.Old_Config_Key)
	return bl, err
}

// DelSystemConfig 删除数据库中的系统参数配置
// @mysql 数据库实例
// @configKey 系统参数名称
func DelSystemConfig(mysql *sql.MySqlDb, configKey string) (bool, error) {
	sqlstr := "delete from system_config where config_key=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, configKey)
	return bl, err
}

// GetModuleInfoList 获取所有功能模块列表
// @mysql 数据库实例
func GetModuleInfoList(mysql *sql.MySqlDb) ([]models.ModuleInfo, error) {
	sqlstr := "SELECT a.*,b.module_name as parent_module_name FROM module_info a left join module_info b on a.parent_module_id = b.module_id"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.ModuleInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ModuleInfo, error) {
		var mo models.ModuleInfo
		mo.Module_Id, err = strconv.Atoi(kr["module_id"])
		if err != nil {
			return *new(models.ModuleInfo), err
		}
		mo.Module_Name = kr["module_name"]
		mo.Parent_Module_Id, _ = strconv.Atoi(kr["parent_module_id"])
		mo.Parent_Module_Name = kr["parent_module_name"]
		mo.Module_Addtime = kr["module_addtime"]
		mo.Module_Icon = kr["module_icon"]
		mo.Module_Info = kr["module_info"]
		mo.Module_Operation = kr["module_operation"]
		operations := make([]string, 0)
		for _, v := range mo.Module_Operation {
			operations = append(operations, string(v))
		}
		mo.Module_Operations = operations
		mo.Module_Url = kr["module_url"]
		if len(kr["sort_num"]) > 0 {
			mo.Sort_Num, err = strconv.Atoi(kr["sort_num"])
			if err != nil {
				return *new(models.ModuleInfo), err
			}
		}
		return mo, nil
	})
	if err != nil {
		return *new([]models.ModuleInfo), err
	}
	return data, nil
}

// AddModuleInfo添加模块功能
// @mysql数据库连接实例
// @info 功能模块实例
func AddModuleInfo(mysql *sql.MySqlDb, info models.ModuleInfo) (bool, error) {
	sqlstr := "insert into module_info(module_name,parent_module_id,module_addtime,module_icon,module_info,module_operation,module_url)values(?,?,?,?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Module_Name, info.Parent_Module_Id, info.Module_Addtime, info.Module_Icon, info.Module_Info, info.Module_Operation, info.Module_Url)
	return bl, err
}

// UpdateModuleInfo更新模块功能
// @mysql数据库连接实例
// @info 功能模块实例
func UpdateModuleInfo(mysql *sql.MySqlDb, info models.ModuleInfo) (bool, error) {
	sqlstr := "update module_info set module_name=?,parent_module_id=?,module_icon=?,module_info=?,module_operation=?,module_url=? where module_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Module_Name, info.Parent_Module_Id, info.Module_Icon, info.Module_Info, info.Module_Operation, info.Module_Url, info.Module_Id)
	return bl, err
}

// DelModuleInfo 删除功能模块
// @mysql数据库连接实例
// @moduleId 功能模块编号
func DelModuleInfo(mysql *sql.MySqlDb, moduleId int) (bool, error) {
	sqlstr := "delete from module_info where module_id=? or parent_module_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, moduleId, moduleId)
	return bl, err
}

// GetModuleFuncs 获取所有模块操作功能列表
// @mysql 数据库实例
func GetModuleFuncs(mysql *sql.MySqlDb) ([]models.ModuleFunc, error) {
	sqlstr := "SELECT a.*,b.module_name as parent_module_name FROM module_func a left join module_info b on a.module_id = b.module_id"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.ModuleFunc), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ModuleFunc, error) {
		var mo models.ModuleFunc
		mo.Func_Id, err = strconv.Atoi(kr["func_id"])
		if err != nil {
			return *new(models.ModuleFunc), err
		}
		mo.Func_Url = kr["func_url"]
		mo.Module_Id, err = strconv.Atoi(kr["module_id"])
		if err != nil {
			return *new(models.ModuleFunc), err
		}
		mo.Module_Name = kr["module_name"]
		mo.Module_Operation = kr["module_operation"]
		mo.Add_Time = kr["add_time"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.ModuleFunc), err
	}
	return data, nil
}

// AddModuleFunc添加模块操作功能
// @mysql数据库连接实例
// @info 功能模块实例
func AddModuleFunc(mysql *sql.MySqlDb, info models.ModuleFunc) (bool, error) {
	sqlstr := "insert into module_func(func_url,module_id,module_operation,add_time)values(?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Func_Url, info.Module_Id, info.Module_Operation, info.Add_Time)
	return bl, err
}

// UpdateModuleFunc更新模块操作功能
// @mysql数据库连接实例
// @info 功能模块实例
func UpdateModuleFunc(mysql *sql.MySqlDb, info models.ModuleFunc) (bool, error) {
	sqlstr := "update module_func set func_url=?,module_id=?,module_operation=? where func_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Func_Url, info.Module_Id, info.Module_Operation, info.Module_Id)
	return bl, err
}

// DelModuleFunc 删除模块的操作功能
// @mysql数据库连接实例
// @moduleId 功能模块编号
func DelModuleFunc(mysql *sql.MySqlDb, funcId int) (bool, error) {
	sqlstr := "delete from module_func where func_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, funcId)
	return bl, err
}

func GetRoleInfoList(mysql *sql.MySqlDb) ([]models.RoleInfo, error) {
	sqlstr := "select a.*,b.role_name as parent_role_name from role_info a left join role_info b on a.parent_id =b.role_id"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.RoleInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.RoleInfo, error) {
		var mo models.RoleInfo
		mo.Role_Id, err = strconv.Atoi(kr["role_id"])
		if err != nil {
			return *new(models.RoleInfo), err
		}
		mo.Role_Name = kr["role_name"]
		mo.Parent_Id, err = strconv.Atoi(kr["parent_id"])
		if err != nil {
			return *new(models.RoleInfo), err
		}
		mo.Parent_Role_Name = kr["parent_role_name"]
		mo.Insert_Time = kr["insert_time"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.RoleInfo), err
	}
	return data, nil
}

func AddRoleInfo(mysql *sql.MySqlDb, info models.RoleInfo) (bool, error) {
	sqlstr := "insert into role_info(role_name,parent_id,insert_time)values(?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Role_Name, info.Parent_Id, info.Insert_Time)
	return bl, err
}

func UpdateRoleInfo(mysql *sql.MySqlDb, info models.RoleInfo) (bool, error) {
	sqlstr := "update role_info set role_name=?,parent_id=? where role_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Role_Name, info.Parent_Id, info.Role_Id)
	return bl, err
}

func DelRoleInfo(mysql *sql.MySqlDb, roleId int) (bool, error) {
	sqlstr := "delete from role_info where role_id=? or parent_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, roleId, roleId)
	return bl, err
}
func GetRoleModules(mysql *sql.MySqlDb, roleId int) ([]models.ModuleInfo, error) {
	sqlstr := "SELECT module_id,module_operation FROM role_module where role_id=?;"
	query, err := mysql.QuerySql(sqlstr, roleId)
	if err != nil {
		return *new([]models.ModuleInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ModuleInfo, error) {
		var mo models.ModuleInfo
		mo.Module_Id, err = strconv.Atoi(kr["module_id"])
		if err != nil {
			return *new(models.ModuleInfo), err
		}
		mo.Module_Operation = kr["module_operation"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.ModuleInfo), err
	}
	return data, nil
}

func UpdateRoleModule(mysql *sql.MySqlDb, mo models.RoleModule) (bool, error) {
	//删除角色功能信息
	//开始事务
	bl, err := mysql.BeginTransact()
	if err != nil {
		return bl, err
	}
	sqlstr := "delete from role_module where role_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, mo.RoleId)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	for _, v := range mo.Modules {
		sqlstr := "insert into role_module(role_id,module_id,module_operation)values(?,?,?)"
		bl, err := mysql.BeginExecSql(sqlstr, mo.RoleId, v.Module_Id, v.Module_Operation)
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return bl, err
		}
	}
	//提交事务
	mysql.CommitTransact()
	return true, nil
}

func GetUserInfoList(mysql *sql.MySqlDb) ([]models.UserInfo, error) {
	sqlstr := "select a.id,a.user_id,a.user_name,a.user_password,a.insert_time,GROUP_CONCAT(b.role_id)as role_ids,GROUP_CONCAT(c.role_name)as role_names from user_info a,user_role b, role_info c where a.id = b.user_id and b.role_id=c.role_id GROUP BY a.id "
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.UserInfo), err
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
		mo.RoleId = kr["role_ids"]
		mo.RoleName = kr["role_names"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.UserInfo), err
	}
	return data, nil
}

func AddUserInfo(mysql *sql.MySqlDb, info models.UserInfo) (bool, error) {
	bl, err := mysql.BeginTransact()
	if err != nil {
		return bl, err
	}
	sqlstr := "insert into user_info(user_id,user_name,user_password,insert_time)values(?,?,?,?);"
	bl, err = mysql.BeginExecSql(sqlstr, info.User_Id, info.User_Name, info.User_Password, info.Insert_Time)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	sqlstr = sql.GetIDENTITYSql()
	query, err := mysql.BeginQueryScalarSql(sqlstr)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	userId, err := sql.MapConvertModel(query, func(kr map[string]string) (int, error) {
		return strconv.Atoi(kr["ID"])
	})
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	if userId == 0 {
		//事务回滚
		mysql.RollbackTransact()
		return false, errors.New("新增用户失败")
	}
	rids := strings.Split(info.RoleId, ",")
	for _, v := range rids {
		sqlstr := "insert into user_role(role_id,user_id,insert_time)values(?,?,?)"
		rid, err := strconv.Atoi(v)
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return false, err
		}
		bl, err = mysql.BeginExecSql(sqlstr, rid, userId, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return bl, err
		}
	}
	mysql.CommitTransact()
	return true, nil
}

func UpdateUserInfo(mysql *sql.MySqlDb, info models.UserInfo) (bool, error) {
	bl, err := mysql.BeginTransact()
	if err != nil {
		return bl, err
	}
	sqlstr := "update user_info set user_id=?,user_name=?,user_password=? where id=?"
	bl, err = mysql.BeginExecSql(sqlstr, info.User_Id, info.User_Name, info.User_Password, info.Id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	//删除原有用户角色表信息
	sqlstr = "delete from user_role where user_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, info.Id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	//添加新的用户角色表信息
	rids := strings.Split(info.RoleId, ",")
	for _, v := range rids {
		sqlstr := "insert into user_role(role_id,user_id,insert_time)values(?,?,?)"
		rid, err := strconv.Atoi(v)
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return false, err
		}
		bl, err = mysql.BeginExecSql(sqlstr, rid, info.Id, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return bl, err
		}
	}
	mysql.CommitTransact()
	return true, nil
}

func DelUserInfo(mysql *sql.MySqlDb, userId int) (bool, error) {
	//删除用户信息
	bl, err := mysql.BeginTransact()
	if err != nil {
		return bl, err
	}
	sqlstr := "delete from user_info where id=?"
	bl, err = mysql.BeginExecSql(sqlstr, userId)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	//删除用户角色表信息
	sqlstr = "delete from user_role where user_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, userId)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	mysql.CommitTransact()
	return true, nil
}
