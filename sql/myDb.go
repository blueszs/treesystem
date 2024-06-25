package sql

import (
	"context"
	db "database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	pcomm "treesystem/common"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/microsoft/go-mssqldb"
)

// OpenDb 根据Init初始化后打开数据库
// @sqlC 数据库连接参数对象
func (d *MySqlDb) OpenDb() error {
	switch d.conf.Dbtype {
	case "mysql":
		dbDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", d.conf.Userid, d.conf.Password, d.conf.Host, d.conf.Port, d.conf.DbName, d.conf.Charset)
	case "sqlserver":
		dbDSN = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", d.conf.Host, d.conf.Userid, d.conf.Password, d.conf.Port, d.conf.DbName)
	default:
		dbDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", d.conf.Userid, d.conf.Password, d.conf.Host, d.conf.Port, d.conf.DbName, d.conf.Charset)
	}
	// 打开连接
	myDb, MysqlDbErr := db.Open(d.conf.Dbtype, dbDSN)
	// 最大连接周期，单位分钟
	myDb.SetConnMaxLifetime(time.Minute * time.Duration(d.conf.ConnMaxLifetime))
	// 最大连接数
	myDb.SetMaxOpenConns(d.conf.MaxOpenConns)
	// 闲置连接数
	myDb.SetMaxIdleConns(d.conf.MaxIdleConns)
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		return MysqlDbErr
	}
	_, MysqlDbErr = CheckSqlAlive(myDb)
	if MysqlDbErr != nil {
		log.Fatal(MysqlDbErr.Error())
		return MysqlDbErr
	}
	d.sqlDB = myDb
	return nil
}

// QuerySql 查询数据,取所有字段,不采用结构体
// @sqlstate sql查询语句
// @params sql查询语句条件
func (d *MySqlDb) QuerySql(sqlstate string, params ...interface{}) (map[int]map[string]string, error) {
	//检查数据库是否处于活动状态
	ctx, err := CheckSqlAlive(d.sqlDB)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//查询数据，取所有字段
	rows2, err := d.sqlDB.QueryContext(ctx, sqlstate, params...)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//返回所有列
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Printf("get columns is error : %v\n", err)
		return nil, err
	}
	//一行所有列的值，用[]byte表示
	vales := make([][]byte, len(cols))
	//一行填充数据
	scans := make([]interface{}, len(cols))
	//scans引用vales，把数据填充到[]byte里
	for k := range vales {
		scans[k] = &vales[k]
	}
	//关闭rows释放持有的数据库链接
	defer CloseSqlRows(rows2)
	i := 0
	result := make(map[int]map[string]string)
	for rows2.Next() {
		//填充数据
		err := rows2.Scan(scans...)
		if err != nil {
			fmt.Printf("fill data row,error: %v\n", err)
			return nil, err
		}
		//每行数据
		row := make(map[string]string)
		//把vales中的数据复制到row中
		for k, v := range vales {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = pcomm.ConvertUnicodeToChar(string(v))
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result, nil
}

func (d *MySqlDb) QueryScalarSql(sqlstate string, params ...interface{}) (map[string]string, error) {
	//检查数据库是否处于活动状态
	ctx, err := CheckSqlAlive(d.sqlDB)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//查询数据，取所有字段
	rows2, err := d.sqlDB.QueryContext(ctx, sqlstate, params...)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//返回所有列
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Printf("get columns is error : %v\n", err)
		return nil, err
	}
	//一行所有列的值，用[]byte表示
	vales := make([][]byte, len(cols))
	//一行填充数据
	scans := make([]interface{}, len(cols))
	//scans引用vales，把数据填充到[]byte里
	for k := range vales {
		scans[k] = &vales[k]
	}
	//关闭rows释放持有的数据库链接
	defer CloseSqlRows(rows2)
	result := make(map[string]string)
	for rows2.Next() {
		//填充数据
		err := rows2.Scan(scans...)
		if err != nil {
			fmt.Printf("fill data row,error: %v\n", err)
			return nil, err
		}
		//每行数据
		//把vales中的数据复制到row中
		for k, v := range vales {
			key := cols[k]
			//这里把[]byte数据转成string
			result[key] = pcomm.ConvertUnicodeToChar(string(v))
		}
		goto ret
	}
ret:
	return result, nil
}

func (d *MySqlDb) QueryModel(f func(map[string]string) any, sqlstate string, params ...any) (any, error) {
	//检查数据库是否处于活动状态
	ctx, err := CheckSqlAlive(d.sqlDB)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return *new(any), err
	}
	//查询数据，取所有字段
	rows2, err := d.sqlDB.QueryContext(ctx, sqlstate, params...)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return *new(any), err
	}
	//返回所有列
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Printf("get columns is error : %v\n", err)
		return *new(any), err
	}
	//一行所有列的值，用[]byte表示
	vales := make([][]byte, len(cols))
	//一行填充数据
	scans := make([]interface{}, len(cols))
	//scans引用vales，把数据填充到[]byte里
	for k := range vales {
		scans[k] = &vales[k]
	}
	//关闭rows释放持有的数据库链接
	defer CloseSqlRows(rows2)
	modes := make([]any, 0)
	for rows2.Next() {
		//填充数据
		err := rows2.Scan(scans...)
		if err != nil {
			fmt.Printf("fill data row,error: %v\n", err)
			return nil, err
		}
		//每行数据
		row := make(map[string]string)
		//把vales中的数据复制到row中
		for k, v := range vales {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = pcomm.ConvertUnicodeToChar(string(v))
		}
		//放入结果集
		//result[i] = row
		rt := f(row)
		modes = append(modes, rt)
	}
	return modes, nil
}

// ExecSqlByTransact 在sql的事务中执行sql语句返回是否执行成功，错误时返回错误
// @sql sql语句
// @params sql语句参数
func (d *MySqlDb) ExecSqlByTransact(sql string, params ...interface{}) (bool, error) {
	//检查数据库是否处于活动状态
	ctx, err := CheckSqlAlive(d.sqlDB)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return false, err
	}
	//开始事务
	tx, err := d.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("begin transactions error:%v\n", err)
		return false, err
	}
	//执行sql
	ret, err := tx.ExecContext(ctx, sql, params...)
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec sql, error:%v\n", err)
		return false, err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("exec sql, error:%v\n", err)
		_ = tx.Rollback()
		return false, err
	}
	//提交事务
	err = tx.Commit()
	if err != nil {
		fmt.Printf("commit transactions, error:%v\n", err)
		return false, err
	}
	fmt.Printf("exec sql success, change rows： %d.\n", rows)
	return true, nil
}

func (d *MySqlDb) BeginTransact() (bool, error) {
	//检查数据库是否处于活动状态
	ctx, err := CheckSqlAlive(d.sqlDB)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return false, err
	}
	d.ctx = &ctx
	//开始事务
	tx, err := d.sqlDB.BeginTx(*d.ctx, nil)
	if err != nil {
		fmt.Printf("begin transactions error:%v\n", err)
		return false, err
	}
	d.Tx = tx
	return true, nil
}

func (d *MySqlDb) BeginExecSql(sql string, params ...interface{}) (bool, error) {
	//执行sql
	ret, err := d.Tx.ExecContext(*d.ctx, sql, params...)
	if err != nil {
		fmt.Printf("exec sql, error:%v\n", err)
		return false, err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("exec sql, error:%v\n", err)
		return false, err
	}
	return rows > 0, nil
}

func (d *MySqlDb) BeginQuerySql(sql string, params ...interface{}) (map[int]map[string]string, error) {
	//执行sql
	//查询数据，取所有字段
	rows2, err := d.Tx.QueryContext(*d.ctx, sql, params...)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//返回所有列
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Printf("get columns is error : %v\n", err)
		return nil, err
	}
	//一行所有列的值，用[]byte表示
	vales := make([][]byte, len(cols))
	//一行填充数据
	scans := make([]interface{}, len(cols))
	//scans引用vales，把数据填充到[]byte里
	for k := range vales {
		scans[k] = &vales[k]
	}
	//关闭rows释放持有的数据库链接
	defer CloseSqlRows(rows2)
	result := make(map[int]map[string]string)
	i := 0
	for rows2.Next() {
		//填充数据
		err := rows2.Scan(scans...)
		if err != nil {
			fmt.Printf("fill data row,error: %v\n", err)
			return nil, err
		}
		//每行数据
		row := make(map[string]string)
		//把vales中的数据复制到row中
		for k, v := range vales {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = pcomm.ConvertUnicodeToChar(string(v))
		}
		//放入结果集
		result[i] = row
	}
	return result, nil
}

func (d *MySqlDb) BeginQueryScalarSql(sql string, params ...interface{}) (map[string]string, error) {
	//执行sql
	//查询数据，取所有字段
	rows2, err := d.Tx.QueryContext(*d.ctx, sql, params...)
	if err != nil {
		fmt.Printf("sql query error: %v\n", err)
		return nil, err
	}
	//返回所有列
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Printf("get columns is error : %v\n", err)
		return nil, err
	}
	//一行所有列的值，用[]byte表示
	vales := make([][]byte, len(cols))
	//一行填充数据
	scans := make([]interface{}, len(cols))
	//scans引用vales，把数据填充到[]byte里
	for k := range vales {
		scans[k] = &vales[k]
	}
	//关闭rows释放持有的数据库链接
	defer CloseSqlRows(rows2)
	result := make(map[string]string)
	for rows2.Next() {
		//填充数据
		err := rows2.Scan(scans...)
		if err != nil {
			fmt.Printf("fill data row,error: %v\n", err)
			return nil, err
		}
		//每行数据
		//把vales中的数据复制到row中
		for k, v := range vales {
			key := cols[k]
			//这里把[]byte数据转成string
			result[key] = pcomm.ConvertUnicodeToChar(string(v))
		}
		goto ret
	}
ret:
	return result, nil
}

func (d *MySqlDb) RollbackTransact() error {
	err := d.Tx.Rollback()
	if err == nil {
		d.Tx = nil
		d.ctx = nil
	}
	return err
}

func (d *MySqlDb) CommitTransact() error {
	err := d.Tx.Commit()
	if err == nil {
		d.Tx = nil
		d.ctx = nil
	}
	return err
}

// CloseSqlConn 关闭sql的连接
func (d *MySqlDb) CloseSqlConn() {
	if d != nil && d.sqlDB != nil {
		err := d.sqlDB.Close()
		if err != nil {
			fmt.Printf("mysql client conn close failed:%v\n", err)
		} else {
			fmt.Println("mysql client conn close success!")
		}
	} else {
		fmt.Println("mysql conn object is null or undefined!")
	}
}

// CheckSqlAlive 检查数据库是否处于活动状态。
// @r 数据库对象
func CheckSqlAlive(r *db.DB) (context.Context, error) {
	ctx := context.Background()
	// Check if database is alive.
	err := r.PingContext(ctx)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

// CloseSqlRows 关闭sql.Rows的连接
func CloseSqlRows(r *db.Rows) error {
	if r != nil {
		err := r.Close()
		if err != nil {
			fmt.Printf("Rows close failed:%v\n", err)
			return err
		}
		return nil
	}
	return errors.New("rows为空，无法关闭。")
}
