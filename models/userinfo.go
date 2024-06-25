package models

type UserInfo struct {
	Id            int    `json:"id"`
	User_Id       string `json:"user_id"`
	User_Name     string `json:"user_name"`
	User_Password string `json:"user_password"`
	Insert_Time   string `json:"insert_time"`
	Token         string `json:"token"`
	RoleId        string `json:"role_id"`
	RoleName      string `json:"role_name"`
}

type RoleModule struct {
	RoleId        int          `json:"role_id"`
	RoleName      string       `json:"role_name"`
	Operations    string       `json:"operations"`
	ModuleAddress string       `json:"module_address"`
	Modules       []ModuleInfo `json:"modules"`
}

type RoleInfo struct {
	Role_Id          int    `json:"role_id"`
	Role_Name        string `json:"role_name"`
	Parent_Id        int    `json:"parent_id"`
	Parent_Role_Name string `json:"parent_role_name"`
	Insert_Time      string `json:"insert_time"`
}

type UpdateRoleModule struct {
	RoleId  int              `json:"role_id"`
	Modules []RoleModuleInfo `json:"modules"`
}

type RoleModuleInfo struct {
	Module_Id        int    `json:"module_id"`
	Module_Operation string `json:"module_operation"`
}
