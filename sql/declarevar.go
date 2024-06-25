package sql

import (
	"context"
	db "database/sql"
)

type MySqlDb struct {
	sqlDB *db.DB
	conf  *SqlConf
	Tx    *db.Tx
	ctx   *context.Context
}

type SqlConf struct {
	Dbtype          string `json:"dbtype"`
	Userid          string `json:"userid"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	DbName          string `json:"dbname"`
	Charset         string `json:"charset"`
	Pooling         int32  `json:"pooling"`
	MultiStatements string `json:"multistatements"`
	Sslmode         string `json:"sslmode"`
	ConnMaxLifetime int    `json:"connmaxlifetime"`
	MaxOpenConns    int    `json:"maxopenconns"`
	MaxIdleConns    int    `json:"maxidleconns"`
}

var dbDSN string
