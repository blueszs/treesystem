package dal

import (
	"strconv"
	"time"
	"treesystem/models"
	"treesystem/sql"
)

// GetListenClassList 获取所有课程类型信息
// @mysql 数据库连接实例
func GetListenClassList(mysql *sql.MySqlDb) ([]models.ListenClass, error) {
	sqlstr := "SELECT a.*  FROM listen_class a;"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.ListenClass), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ListenClass, error) {
		var mo models.ListenClass
		mo.Class_Id, err = strconv.Atoi(kr["class_id"])
		if err != nil {
			return *new(models.ListenClass), err
		}
		mo.Class_Name = kr["class_name"]
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.ListenClass), err
		}
		switch mo.Flag {
		case 0:
			mo.Flag_Info = "禁用"
		case 1:
			mo.Flag_Info = "维护"
		case 2:
			mo.Flag_Info = "正常"
		}
		return mo, nil
	})
	if err != nil {
		return *new([]models.ListenClass), err
	}
	return data, nil
}

// AddListenClass 添加课程类型
// @mysql数据库连接实例
// @info 课程类型实例
func AddListenClass(mysql *sql.MySqlDb, info models.ListenClass) (bool, error) {
	sqlstr := "insert into listen_class(class_name,flag) values(?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Class_Name, info.Flag)
	return bl, err
}

// UpdateListenClass 修改课程类型
// @mysql数据库连接实例
// @info 课程类型实例
func UpdateListenClass(mysql *sql.MySqlDb, info models.ListenClass) (bool, error) {
	sqlstr := "update from listen_class set class_name=?,flag=? where class_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Class_Name, info.Flag, info.Class_Id)
	return bl, err
}

// DelListenClass 删除课程类型
// @mysql数据库连接实例
// @memberId 课程类型编号
func DelListenClass(mysql *sql.MySqlDb, memberId int) (bool, error) {
	sqlstr := "delete from listen_class where member_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, memberId)
	return bl, err
}

// GetListenSeriesList 获取所有课程系列信息
// @mysql 数据库连接实例
func GetListenSeriesList(mysql *sql.MySqlDb) ([]models.ListenSeries, error) {
	sqlstr := "SELECT a.*,b.class_name  FROM listen_series a left join listen_class b on a.class_id=b.class_id;"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.ListenSeries), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ListenSeries, error) {
		var mo models.ListenSeries
		mo.Series_Id, err = strconv.Atoi(kr["series_id"])
		if err != nil {
			return *new(models.ListenSeries), err
		}
		mo.Series_Name = kr["series_name"]
		mo.Class_Id, err = strconv.Atoi(kr["class_id"])
		if err != nil {
			return *new(models.ListenSeries), err
		}
		mo.Class_Name = kr["class_name"]
		mo.Series_Info = kr["series_info"]
		mo.Series_Img = kr["series_img"]
		mo.Add_Time = kr["add_time"]
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.ListenSeries), err
		}
		switch mo.Flag {
		case 0:
			mo.Flag_Info = "禁用"
		case 1:
			mo.Flag_Info = "维护"
		case 2:
			mo.Flag_Info = "正常"
		}
		return mo, nil
	})
	if err != nil {
		return *new([]models.ListenSeries), err
	}
	return data, nil
}

// AddListenSeries 添加课程系列
// @mysql数据库连接实例
// @info 课程系列实例
func AddListenSeries(mysql *sql.MySqlDb, info models.ListenSeries) (bool, error) {
	sqlstr := "insert into listen_series(series_name,class_id,series_info,series_img,flag,add_time) values(?,?,?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Series_Name, info.Class_Id, info.Series_Info, info.Series_Img, info.Flag, time.Now().Format("2006-01-02 15:04:05"))
	return bl, err
}

// UpdateListenSeries 修改课程系列
// @mysql数据库连接实例
// @info 课程系列实例
func UpdateListenSeries(mysql *sql.MySqlDb, info models.ListenSeries) (bool, error) {
	sqlstr := "update from listen_series set series_name=? class_id=?,series_info=?,series_img=?,flag=? where series_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Series_Name, info.Class_Id, info.Series_Info, info.Series_Img, info.Flag, info.Series_Id)
	return bl, err
}

// DelListenSeries 删除课程系列
// @mysql数据库连接实例
// @seriesId 课程系列编号
func DelListenSeries(mysql *sql.MySqlDb, seriesId int) (bool, error) {
	//开始事务
	bl, err := mysql.BeginTransact()
	if err != nil {
		return bl, err
	}
	// 删除课程系列
	sqlstr := "delete from listen_series where series_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, seriesId)
	if err != nil {
		mysql.RollbackTransact()
		return bl, err
	}
	// 删除课程系列中的课程
	sqlstr = "delete from listen_info where series_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, seriesId)
	if err != nil {
		mysql.RollbackTransact()
		return bl, err
	}
	// 提交事务
	mysql.CommitTransact()
	return bl, err
}

// GetListenInfoList 获取所有课程信息
// @mysql 数据库连接实例
func GetListenInfoList(mysql *sql.MySqlDb) ([]models.ListenInfo, error) {
	sqlstr := "SELECT a.*,b.series_name  FROM listen_info a left join listen_series b on a.series_id=b.series_id;"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.ListenInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.ListenInfo, error) {
		var mo models.ListenInfo
		mo.Listen_Id, err = strconv.Atoi(kr["listen_id"])
		if err != nil {
			return *new(models.ListenInfo), err
		}
		mo.Listen_Name = kr["listen_name"]
		mo.Series_Id, err = strconv.Atoi(kr["series_id"])
		if err != nil {
			return *new(models.ListenInfo), err
		}
		mo.Series_Name = kr["series_name"]
		mo.Listen_Info = kr["listen_info"]
		mo.Listen_Img = kr["listen_img"]
		mo.Listen_Ep, err = strconv.Atoi(kr["listen_ep"])
		if err != nil {
			return *new(models.ListenInfo), err
		}
		mo.Add_Time = kr["add_time"]
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.ListenInfo), err
		}
		switch mo.Flag {
		case 0:
			mo.Flag_Info = "禁用"
		case 1:
			mo.Flag_Info = "维护"
		case 2:
			mo.Flag_Info = "正常"
		}
		return mo, nil
	})
	if err != nil {
		return *new([]models.ListenInfo), err
	}
	return data, nil
}

// AddListenInfo 添加课程
// @mysql数据库连接实例
// @info 课程实例
func AddListenInfo(mysql *sql.MySqlDb, info models.ListenInfo) (bool, error) {
	sqlstr := "insert into listen_info(listen_name,series_id,listen_info,listen_img,listen_ep,flag,add_time) values(?,?,?,?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Listen_Name, info.Series_Id, info.Listen_Info, info.Listen_Img, info.Listen_Ep, info.Flag, time.Now().Format("2006-01-02 15:04:05"))
	return bl, err
}

// UpdateListenInfo 修改课程
// @mysql数据库连接实例
// @info 课程实例
func UpdateListenInfo(mysql *sql.MySqlDb, info models.ListenInfo) (bool, error) {
	sqlstr := "update from listen_info set listen_name=? series_id=?,listen_info=?,listen_img=?,listen_ep=?,flag=? where listen_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Listen_Name, info.Series_Id, info.Listen_Info, info.Listen_Img, info.Listen_Ep, info.Flag, info.Listen_Id)
	return bl, err
}

// DelListenInfo 删除课程
// @mysql数据库连接实例
// @listenId 课程编号
func DelListenInfo(mysql *sql.MySqlDb, listenId int) (bool, error) {
	sqlstr := "delete from listen_info where listen_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, listenId)
	return bl, err
}
