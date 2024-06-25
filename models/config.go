package models

import "treesystem/sql"

type Config struct {
	Database sql.SqlConf `json:"database"`
	Web      WebConfig   `json:"web"`
	System   SysConfig   `json:"system"`
	Cache    []CacheConf `json:"cache"`
}

type WebConfig struct {
	Port         int `json:"port"`
	ReadTimeout  int `json:"readtimeout"`
	WriteTimeout int `json:"writetimeout"`
}

type SysConfig struct {
	PwdKey   string `json:"pwdkey"`
	TokenKey string `json:"tokenkey"`
}

type ModuleInfo struct {
	Module_Id          int      `json:"module_id"`
	Module_Name        string   `json:"module_name"`
	Parent_Module_Id   int      `json:"parent_module_id"`
	Parent_Module_Name string   `json:"parent_module_name"`
	Module_Url         string   `json:"module_url"`
	Module_Info        string   `json:"module_info"`
	Module_Icon        string   `json:"module_icon"`
	Module_Addtime     string   `json:"module_addtime"`
	Module_Operation   string   `json:"module_operation"`
	Module_Operations  []string `json:"module_operations"`
	Sort_Num           int      `json:"sort_num"`
}

type ModuleFunc struct {
	Func_Id          int    `json:"func_id"`
	Func_Url         string `json:"func_url"`
	Module_Id        int    `json:"module_id"`
	Module_Name      string `json:"module_name"`
	Module_Operation string `json:"module_operation"`
	Add_Time         string `json:"add_time"`
}

type CacheConf struct {
	IpAddr   string `json:"ip"`
	IpPort   string `json:"port"`
	Password string `json:"password"`
	DbIndex  int    `json:"dbindex"`
}
