package dal

import (
	"strconv"
	"time"
	"treesystem/models"
	"treesystem/sql"
)

// GetTreeList 根据关键字搜索，获取树种列表
// @mysql 数据库实例结构体
// @skey 搜索关键字
func GetTreeList(mysql *sql.MySqlDb, skey string) ([]models.TreeInfo, error) {
	sqlstr := "select a.tree_id,a.tree_name,a.class_id,a.tree_info,a.tree_main_photo,a.tree_main_photo_thumbnail,b.class_name from tree_info a left join tree_class b on a.class_id=b.class_id where a.flag>1 and a.tree_name like ?"
	query, err := mysql.QuerySql(sqlstr, "%"+skey+"%")
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
		mo.Tree_Main_Photo = kr["tree_main_photo"]
		mo.Tree_Main_Photo_Thumbnail = kr["tree_main_photo_thumbnail"]
		return mo, nil
	})
	if err != nil {
		return *new([]models.TreeInfo), err
	}
	return data, nil
}

func GetTreeInfo(mysql *sql.MySqlDb, treeid int) (models.TreeInfo, error) {
	sqlstr := "select a.*,b.class_name from tree_info a left join tree_class b on a.class_id=b.class_id where a.tree_id=? and a.flag>1;"
	query, err := mysql.QueryScalarSql(sqlstr, treeid)
	if err != nil {
		return *new(models.TreeInfo), err
	}
	info, err := sql.MapConvertModel(query, func(kr map[string]string) (models.TreeInfo, error) {
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
		mo.Tree_Main_Photo = kr["tree_main_photo"]
		mo.Tree_Main_Photo_Thumbnail = kr["tree_main_photo_thumbnail"]
		return mo, nil
	})
	sqlstr = "select * from tree_detail where tree_id=? and flag>1;"
	query1, err1 := mysql.QuerySql(sqlstr, treeid)
	if err1 != nil {
		return *new(models.TreeInfo), err1
	}
	attribute := make([]models.TreeAttribute, 0)
	data, err := sql.ConvertModel(query1, func(kr map[string]string) (models.TreeDetail, error) {
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
		if len(mo.Detail_Dim_Key1) > 0 {
			var treeAttr models.TreeAttribute
			index := 0
			for i, attr := range attribute {
				if attr.Name == mo.Detail_Dim_Key1 {
					treeAttr = attr
					index = i
					break
				}
			}
			var attinfo models.AtrrInfo
			attinfo.Attrdetal = mo.Detail_Dim_Value1
			treeAttr.Values = append(treeAttr.Values, attinfo)
			if len(treeAttr.Name) == 0 {
				treeAttr.Name = mo.Detail_Dim_Key1
				attribute = append(attribute, treeAttr)
			} else {
				attribute[index] = treeAttr
			}
		}
		if len(mo.Detail_Dim_Key2) > 0 {
			var treeAttr models.TreeAttribute
			index := 0
			for i, attr := range attribute {
				if attr.Name == mo.Detail_Dim_Key2 {
					treeAttr = attr
					index = i
					break
				}
			}
			var attinfo models.AtrrInfo
			attinfo.Attrdetal = mo.Detail_Dim_Value2
			treeAttr.Values = append(treeAttr.Values, attinfo)
			if len(treeAttr.Name) == 0 {
				treeAttr.Name = mo.Detail_Dim_Key2
				attribute = append(attribute, treeAttr)
			} else {
				attribute[index] = treeAttr
			}
		}
		if len(mo.Detail_Dim_Key3) > 0 {
			var treeAttr models.TreeAttribute
			index := 0
			for i, attr := range attribute {
				if attr.Name == mo.Detail_Dim_Key3 {
					treeAttr = attr
					index = i
					break
				}
			}
			var attinfo models.AtrrInfo
			attinfo.Attrdetal = mo.Detail_Dim_Value3
			treeAttr.Values = append(treeAttr.Values, attinfo)
			if len(treeAttr.Name) == 0 {
				treeAttr.Name = mo.Detail_Dim_Key3
				attribute = append(attribute, treeAttr)
			} else {
				attribute[index] = treeAttr
			}
		}
		if len(mo.Detail_Dim_Key4) > 0 {
			var treeAttr models.TreeAttribute
			index := 0
			for i, attr := range attribute {
				if attr.Name == mo.Detail_Dim_Key4 {
					treeAttr = attr
					index = i
					break
				}
			}
			var attinfo models.AtrrInfo
			attinfo.Attrdetal = mo.Detail_Dim_Value4
			treeAttr.Values = append(treeAttr.Values, attinfo)
			if len(treeAttr.Name) == 0 {
				treeAttr.Name = mo.Detail_Dim_Key4
				attribute = append(attribute, treeAttr)
			} else {
				attribute[index] = treeAttr
			}
		}
		if len(mo.Detail_Dim_Key5) > 0 {
			var treeAttr models.TreeAttribute
			index := 0
			for i, attr := range attribute {
				if attr.Name == mo.Detail_Dim_Key5 {
					treeAttr = attr
					index = i
					break
				}
			}
			var attinfo models.AtrrInfo
			attinfo.Attrdetal = mo.Detail_Dim_Value5
			treeAttr.Values = append(treeAttr.Values, attinfo)
			if len(treeAttr.Name) == 0 {
				treeAttr.Name = mo.Detail_Dim_Key5
				attribute = append(attribute, treeAttr)
			} else {
				attribute[index] = treeAttr
			}
		}
		return mo, nil
	})
	if err1 != nil {
		return *new(models.TreeInfo), err1
	}
	info.Attribute = attribute
	info.Details = data
	sqlstr = "select * from tree_photo where detail_id in(select detail_id from tree_detail where tree_id=? and flag>1);"
	query1, err1 = mysql.QuerySql(sqlstr, treeid)
	if err1 != nil {
		return *new(models.TreeInfo), err1
	}
	photos, err := sql.ConvertModel(query1, func(kr map[string]string) (string, error) {
		return kr["img_source"], nil
	})
	if err1 != nil {
		return *new(models.TreeInfo), err1
	}
	info.Tree_Photos = photos
	return info, nil
}

func GetConfigValue(mysql *sql.MySqlDb, configKey string) (string, error) {
	sqlstr := "select config_value from system_config where config_key=?;"
	data, err := mysql.QueryScalarSql(sqlstr, configKey)
	if err != nil {
		return "", err
	}
	val, err := sql.MapConvertModel(data, func(kr map[string]string) (string, error) {
		return kr["config_value"], nil
	})
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetMemberInfo(mysql *sql.MySqlDb, openid string) (models.MemberInfo, error) {
	sqlstr := "select * from member_info where member_wechat_id=?;"
	data, err := mysql.QueryScalarSql(sqlstr, openid)
	if err != nil {
		return *new(models.MemberInfo), err
	}
	member, err := sql.MapConvertModel(data, func(kr map[string]string) (models.MemberInfo, error) {
		var mo models.MemberInfo
		mo.Member_Id, err = strconv.Atoi(kr["member_id"])
		if err != nil {
			return *new(models.MemberInfo), err
		}
		mo.Member_Name = kr["member_name"]
		mo.Member_Nick = kr["member_nick"]
		mo.Member_Tel = kr["member_tel"]
		mo.Member_WeChat_Id = kr["member_wechat_id"]
		mo.Member_Photo = kr["member_photo"]
		mo.Create_Time = kr["create_time"]
		return mo, nil
	})
	if err != nil {
		return *new(models.MemberInfo), err
	}
	return member, nil
}

func UpdateMemberInfo(mysql *sql.MySqlDb, member models.MemberInfo) (bool, error) {
	sqlstr := "select member_id from member_info where member_wechat_id=?;"
	data, err := mysql.QueryScalarSql(sqlstr, member.Member_WeChat_Id)
	if err != nil {
		return false, err
	}
	memberId, err := sql.MapConvertModel(data, func(kr map[string]string) (int, error) {
		return strconv.Atoi(kr["member_id"])
	})
	if err != nil {
		return false, err
	}
	if memberId > 0 {
		sqlstr = "update member_info set member_name=?,member_nick=?,member_tel=?,member_photo=? where member_wechat_id=?"
		_, err = mysql.ExecSqlByTransact(sqlstr, member.Member_Name, member.Member_Nick, member.Member_Tel, member.Member_Photo, member.Member_WeChat_Id)
		if err != nil {
			return false, err
		}
	} else {
		sqlstr = "insert into member_info(member_name,member_nick,member_tel,member_photo,member_wechat_id,create_time)values(@member_name,@member_nick,@member_tel,@member_photo,@member_wechat_id,@create_time)"
		_, err = mysql.ExecSqlByTransact(sqlstr, member.Member_Name, member.Member_Nick, member.Member_Tel, member.Member_Photo, member.Member_WeChat_Id, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
