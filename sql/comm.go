package sql

import db "database/sql"

// Init 初始化mysql连接字符串
// @sqlC 数据库连接参数对象
func Init(sqlC SqlConf) *MySqlDb {
	mysDb := &db.DB{}
	return &MySqlDb{sqlDB: mysDb, conf: &sqlC}
}

// ConvertModel 将字典类型转化为对象数组
// @data 从数据库查询出的字典
// @f 将字典赋值到结构体的方法
func ConvertModel[T any](data map[int]map[string]string, f func(map[string]string) (T, error)) ([]T, error) {
	result := make([]T, len(data))
	for k := range data {
		rt, err := f(data[k])
		if err != nil {
			return *new([]T), err
		}
		result[k] = rt
	}
	return result, nil
}

func MapConvertModel[T any](data map[string]string, f func(map[string]string) (T, error)) (T, error) {
	name, err := f(data)
	if err != nil {
		return *new(T), err
	}
	return name, err
}
func GetIDENTITYSql() string {
	return "SELECT LAST_INSERT_ID() AS ID FROM DUAL;"
}
