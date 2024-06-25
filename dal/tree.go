package dal

import (
	"strconv"
	"time"
	"treesystem/models"
	"treesystem/sql"
)

func GetTreeClassList(mysql *sql.MySqlDb) ([]models.TreeClass, error) {
	sqlstr := "select a.*,b.class_name as parent_class_name from tree_class a left join tree_class b on a.parent_class_id =b.class_id"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.TreeClass), err
	}
	list, err := sql.ConvertModel(query, func(kr map[string]string) (models.TreeClass, error) {
		var mo models.TreeClass
		mo.Class_Id, err = strconv.Atoi(kr["class_id"])
		if err != nil {
			return *new(models.TreeClass), err
		}
		mo.Class_Name = kr["class_name"]
		mo.Parent_Class_Id, err = strconv.Atoi(kr["parent_class_id"])
		mo.Parent_Class_Name = kr["parent_class_name"]
		mo.Insert_Time = kr["insert_time"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.TreeClass), err
	}
	return list, nil
}

func AddTreeClass(mysql *sql.MySqlDb, info models.TreeClass) (bool, error) {
	sqlstr := "insert into tree_class(class_name,parent_class_id,insert_time)values(?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Class_Name, info.Parent_Class_Id, time.Now().Format("2006-01-02 15:04:05"))
	return bl, err
}

func UpdateTreeClass(mysql *sql.MySqlDb, info models.TreeClass) (bool, error) {
	sqlstr := "update tree_class set class_name=?,parent_class_id=? where class_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Class_Name, info.Parent_Class_Id, info.Class_Id)
	return bl, err
}

func DelTreeClass(mysql *sql.MySqlDb, id int) (bool, error) {
	sqlstr := "delete from tree_class where class_id=@class_id or parent_class_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, id)
	return bl, err
}

func GetTreePlantList(mysql *sql.MySqlDb) ([]models.TreePlant, error) {
	sqlstr := "select a.* from tree_plant a"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.TreePlant), err
	}
	list, err := sql.ConvertModel(query, func(kr map[string]string) (models.TreePlant, error) {
		var mo models.TreePlant
		mo.Plant_Id, err = strconv.Atoi(kr["plant_id"])
		if err != nil {
			return *new(models.TreePlant), err
		}
		mo.Plant_Name = kr["plant_name"]
		mo.Plant_Info = kr["plant_info"]
		mo.Plant_Address = kr["plant_address"]
		mo.Plant_Contacts = kr["plant_contacts"]
		mo.Plant_Tel = kr["plant_tel"]
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.TreePlant), err
		}
		mo.Insert_Time = kr["insert_time"]
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
		return *new([]models.TreePlant), err
	}
	return list, nil
}

func AddTreePlant(mysql *sql.MySqlDb, info models.TreePlant) (bool, error) {
	sqlstr := "insert into tree_plant(plant_name,plant_info,plant_address,plant_contacts,plant_tel,flag,insert_time)values(?,?,?,?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Plant_Name, info.Plant_Info, info.Plant_Address, info.Plant_Contacts, info.Plant_Tel, info.Flag, time.Now().Format("2006-01-02 15:04:05"))
	return bl, err
}
func UpdateTreePlant(mysql *sql.MySqlDb, info models.TreePlant) (bool, error) {
	sqlstr := "update tree_plant set plant_name=?,plant_info=?,plant_address=?,plant_contacts=?,plant_tel=?,flag=? where plant_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Plant_Name, info.Plant_Info, info.Plant_Address, info.Plant_Contacts, info.Plant_Tel, info.Flag, info.Plant_Id)
	return bl, err
}

func DelTreePlant(mysql *sql.MySqlDb, id int) (bool, error) {
	sqlstr := "delete from tree_plant where plant_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, id)
	return bl, err
}

func GetTreeInfoList(mysql *sql.MySqlDb) ([]models.TreeInfo, error) {
	sqlstr := "select a.*,b.class_name from tree_info a left join tree_class b on a.class_id=b.class_id"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.TreeInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.TreeInfo, error) {
		var mo models.TreeInfo
		mo.Tree_Id, err = strconv.Atoi(kr["tree_id"])
		if err != nil {
			return *new(models.TreeInfo), err
		}
		mo.Tree_Name = kr["tree_name"]
		mo.Class_Id, err = strconv.Atoi(kr["class_id"])
		if err != nil {
			return *new(models.TreeInfo), err
		}
		mo.Class_Name = kr["class_name"]
		mo.Tree_Info = kr["tree_info"]
		if len(mo.Tree_Info) > 16 {
			var r = []rune(mo.Tree_Info)
			mo.Short_Tree_Info = string(r[0:16]) + "……"
		}
		mo.Tree_Main_Photo = kr["tree_main_photo"]
		mo.Tree_Main_Photo_Thumbnail = kr["tree_main_photo_thumbnail"]
		mo.Flag, err = strconv.Atoi(kr["flag"])
		mo.Insert_Time = kr["insert_time"]
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
		return *new([]models.TreeInfo), err
	}
	return data, nil
}

func AddTreeInfo(mysql *sql.MySqlDb, info models.TreeInfo) (bool, error) {
	sqlstr := "insert into tree_info(tree_name,tree_info,class_id,tree_main_photo,tree_main_photo_thumbnail,flag,insert_time)values(?,?,?,?,?,?,?);"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Tree_Name, info.Tree_Info, info.Class_Id, info.Tree_Main_Photo, info.Tree_Main_Photo_Thumbnail, info.Flag, time.Now().Format("2006-01-02 15:04:05"))
	return bl, err
}

func UpdateTreeInfo(mysql *sql.MySqlDb, info models.TreeInfo) (bool, error) {
	sqlstr := "update tree_info set tree_name=?,tree_info=?,class_id=?,tree_main_photo=?,tree_main_photo_thumbnail=?,flag=? where tree_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, info.Tree_Name, info.Tree_Info, info.Class_Id, info.Tree_Main_Photo, info.Tree_Main_Photo_Thumbnail, info.Flag, info.Tree_Id)
	return bl, err
}

func DelTreeInfo(mysql *sql.MySqlDb, id int) (bool, error) {
	//删除苗木信息
	sqlstr := "delete from tree_info where tree_id=?"
	mysql.BeginTransact()
	bl, err := mysql.BeginExecSql(sqlstr, id)
	if err != nil {
		mysql.RollbackTransact()
		return bl, err
	}
	//删除苗木图片信息
	sqlstr = "delete from tree_photo where tree_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, id)
	if err != nil {
		mysql.RollbackTransact()
		return bl, err
	}
	//删除苗木详细信息
	sqlstr = "delete from tree_detail where tree_id=@tree_id"
	bl, err = mysql.BeginExecSql(sqlstr, id)
	if err != nil {
		mysql.RollbackTransact()
		return bl, err
	}
	mysql.CommitTransact()
	return bl, nil
}

func GetTreeDetailList(mysql *sql.MySqlDb) ([]models.TreeDetail, error) {
	sqlstr := "select a.*,b.tree_name,c.plant_name from tree_detail a left join tree_info b on a.tree_id=b.tree_id left join tree_plant c on a.plant_id=c.plant_id;"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.TreeDetail), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.TreeDetail, error) {
		var mo models.TreeDetail
		mo.Detail_Id, err = strconv.Atoi(kr["detail_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Tree_Id, err = strconv.Atoi(kr["tree_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Tree_Name = kr["tree_name"]
		mo.Plant_Id, err = strconv.Atoi(kr["plant_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Plant_Name = kr["plant_name"]
		mo.Detail_Dim_Key1 = kr["detail_dim_key1"]
		mo.Detail_Dim_Value1 = kr["detail_dim_value1"]
		mo.Detail_Dim_Key2 = kr["detail_dim_key2"]
		mo.Detail_Dim_Value2 = kr["detail_dim_value2"]
		mo.Detail_Dim_Key3 = kr["detail_dim_key3"]
		mo.Detail_Dim_Value3 = kr["detail_dim_value3"]
		mo.Detail_Dim_Key4 = kr["detail_dim_key4"]
		mo.Detail_Dim_Value4 = kr["detail_dim_value4"]
		mo.Detail_Dim_Key5 = kr["detail_dim_key5"]
		mo.Detail_Dim_Value5 = kr["detail_dim_value5"]
		mo.Tree_Unity = kr["tree_unity"]
		mo.Store_Quantity, err = strconv.Atoi(kr["store_quantity"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Sales_Quantity, err = strconv.Atoi(kr["sales_quantity"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Active_Price, err = strconv.ParseFloat(kr["active_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Prime_Price, err = strconv.ParseFloat(kr["prime_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Real_Price, err = strconv.ParseFloat(kr["real_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Insert_Time = kr["insert_time"]
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
		return *new([]models.TreeDetail), err
	}
	return data, nil
}

func GetTreeDetail(mysql *sql.MySqlDb, id int) (models.TreeDetail, error) {
	sqlstr := "select * from tree_photo where detail_id=?;"
	query, err := mysql.QuerySql(sqlstr, id)
	if err != nil {
		return *new(models.TreeDetail), err
	}
	photos, err := sql.ConvertModel(query, func(kr map[string]string) (models.TreeDetailPhoto, error) {
		var mo models.TreeDetailPhoto
		mo.Id, err = strconv.Atoi(kr["id"])
		if err != nil {
			return *new(models.TreeDetailPhoto), err
		}
		mo.Detail_Id, err = strconv.Atoi(kr["detail_id"])
		if err != nil {
			return *new(models.TreeDetailPhoto), err
		}
		mo.Img_Source = kr["img_source"]
		mo.Img_Thumbnail = kr["img_thumbnail"]
		if en := len(kr["sort_num"]); en > 0 {
			mo.Sort_Num, err = strconv.Atoi(kr["sort_num"])
			if err != nil {
				return *new(models.TreeDetailPhoto), err
			}
		}
		return mo, nil
	})
	if err != nil {
		return *new(models.TreeDetail), err
	}
	sqlstr = "select a.*,b.tree_name,c.plant_name from tree_detail a left join tree_info b on a.tree_id=b.tree_id left join tree_plant c on a.plant_id=c.plant_id where detail_id=?;"
	queryMP, err := mysql.QueryScalarSql(sqlstr, id)
	data, err := sql.MapConvertModel(queryMP, func(kr map[string]string) (models.TreeDetail, error) {
		var mo models.TreeDetail
		mo.Detail_Id, err = strconv.Atoi(kr["detail_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Tree_Id, err = strconv.Atoi(kr["tree_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Tree_Name = kr["tree_name"]
		mo.Plant_Id, err = strconv.Atoi(kr["plant_id"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Plant_Name = kr["plant_name"]
		mo.Detail_Dim_Key1 = kr["detail_dim_key1"]
		mo.Detail_Dim_Value1 = kr["detail_dim_value1"]
		mo.Detail_Dim_Key2 = kr["detail_dim_key2"]
		mo.Detail_Dim_Value2 = kr["detail_dim_value2"]
		mo.Detail_Dim_Key3 = kr["detail_dim_key3"]
		mo.Detail_Dim_Value3 = kr["detail_dim_value3"]
		mo.Detail_Dim_Key4 = kr["detail_dim_key4"]
		mo.Detail_Dim_Value4 = kr["detail_dim_value4"]
		mo.Detail_Dim_Key5 = kr["detail_dim_key5"]
		mo.Detail_Dim_Value5 = kr["detail_dim_value5"]
		mo.Tree_Unity = kr["tree_unity"]
		mo.Store_Quantity, err = strconv.Atoi(kr["store_quantity"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Sales_Quantity, err = strconv.Atoi(kr["sales_quantity"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Active_Price, err = strconv.ParseFloat(kr["active_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Prime_Price, err = strconv.ParseFloat(kr["prime_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Real_Price, err = strconv.ParseFloat(kr["real_price"], 32)
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Flag, err = strconv.Atoi(kr["flag"])
		if err != nil {
			return *new(models.TreeDetail), err
		}
		mo.Insert_Time = kr["insert_time"]
		switch mo.Flag {
		case 0:
			mo.Flag_Info = "禁用"
		case 1:
			mo.Flag_Info = "维护"
		case 2:
			mo.Flag_Info = "正常"
		}
		mo.TreeDetailPhotos = photos
		return mo, nil
	})
	return data, nil
}

func AddTreeDetail(mysql *sql.MySqlDb, info models.TreeDetail) (bool, error) {
	sqlstr := "insert into tree_detail(tree_id,plant_id,detail_dim_key1,detail_dim_value1,detail_dim_key2,detail_dim_value2,detail_dim_key3,detail_dim_value3,detail_dim_key4,detail_dim_value4,detail_dim_key5,detail_dim_value5,tree_unity,store_quantity,sales_quantity,active_price,prime_price,real_price,flag,insert_time)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	mysql.BeginTransact()
	bl, err := mysql.BeginExecSql(sqlstr, info.Tree_Id, info.Plant_Id, info.Detail_Dim_Key1, info.Detail_Dim_Value1, info.Detail_Dim_Key2, info.Detail_Dim_Value2, info.Detail_Dim_Key3, info.Detail_Dim_Value3, info.Detail_Dim_Key4, info.Detail_Dim_Value4, info.Detail_Dim_Key5, info.Detail_Dim_Value5, info.Tree_Unity, info.Store_Quantity, info.Sales_Quantity, info.Active_Price, info.Prime_Price, info.Real_Price, info.Flag, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	sqlstr = sql.GetIDENTITYSql()
	query, err := mysql.BeginQueryScalarSql(sqlstr)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	detail_id, err := sql.MapConvertModel(query, func(kr map[string]string) (int, error) {
		return strconv.Atoi(kr["ID"])
	})
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	for _, item := range info.TreeDetailPhotos {
		sqlstr = "insert into tree_photo(detail_id,img_source)values(?,?);"
		bl, err = mysql.BeginExecSql(sqlstr, detail_id, item.Img_Source)
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return false, err
		}
	}
	mysql.CommitTransact()
	return bl, nil
}

func UpdateTreeDetail(mysql *sql.MySqlDb, info models.TreeDetail) (bool, error) {
	sqlstr := "update tree_detail set tree_id=?,plant_id=?,detail_dim_key1=?,detail_dim_value1=?,detail_dim_key2=?," +
		"detail_dim_value2=?,detail_dim_key3=?,detail_dim_value3=?,detail_dim_key4=?,detail_dim_value4=?," +
		"detail_dim_key5=?,detail_dim_value5=?,tree_unity=?,store_quantity=?,sales_quantity=?,active_price=?," +
		"prime_price=?,real_price=?,flag=? where detail_id=?;"
	mysql.BeginTransact()
	bl, err := mysql.BeginExecSql(sqlstr, info.Tree_Id, info.Plant_Id, info.Detail_Dim_Key1, info.Detail_Dim_Value1, info.Detail_Dim_Key2, info.Detail_Dim_Value2, info.Detail_Dim_Key3, info.Detail_Dim_Value3, info.Detail_Dim_Key4, info.Detail_Dim_Value4, info.Detail_Dim_Key5, info.Detail_Dim_Value5, info.Tree_Unity, info.Store_Quantity, info.Sales_Quantity, info.Active_Price, info.Prime_Price, info.Real_Price, info.Flag, info.Detail_Id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	//删除原照片
	sqlstr = "delete from tree_photo where detail_id=?;"
	bl, err = mysql.BeginExecSql(sqlstr, info.Detail_Id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	for _, item := range info.TreeDetailPhotos {
		sqlstr = "insert into tree_photo(detail_id,img_source)values(?,?);"
		bl, err = mysql.BeginExecSql(sqlstr, info.Detail_Id, item.Img_Source)
		if err != nil {
			//事务回滚
			mysql.RollbackTransact()
			return false, err
		}
	}
	mysql.CommitTransact()
	return bl, nil
}

func DelTreeDetail(mysql *sql.MySqlDb, id int) (bool, error) {
	//删除苗木信息
	sqlstr := "delete from tree_detail where detail_id=?"
	bl, err := mysql.BeginExecSql(sqlstr, id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return bl, err
	}
	//删除苗木图片信息
	sqlstr = "delete from tree_photo where detail_id=?"
	bl, err = mysql.BeginExecSql(sqlstr, id)
	if err != nil {
		//事务回滚
		mysql.RollbackTransact()
		return false, err
	}
	return bl, nil
}
