package dal

import (
	"strconv"
	"treesystem/models"
	"treesystem/sql"
)

// GetMemberInfoList 获取所有会员信息
// @mysql 数据库连接实例
func GetMemberInfoList(mysql *sql.MySqlDb) ([]models.MemberInfo, error) {
	sqlstr := "SELECT a.*  FROM Member_info a;"
	query, err := mysql.QuerySql(sqlstr)
	if err != nil {
		return *new([]models.MemberInfo), err
	}
	data, err := sql.ConvertModel(query, func(kr map[string]string) (models.MemberInfo, error) {
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
		return *new([]models.MemberInfo), err
	}
	return data, nil
}

// DelMemberInfo 删除会员账号
// @mysql数据库连接实例
// @memberId 会员编号
func DelMemberInfo(mysql *sql.MySqlDb, memberId int) (bool, error) {
	sqlstr := "delete from member_info where member_id=?"
	bl, err := mysql.ExecSqlByTransact(sqlstr, memberId)
	return bl, err
}
