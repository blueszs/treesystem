package models

// User stores information of a users
type User struct {
	Name  string
	Email string
	Intro string
}

type UserNavbar struct {
	RoleId          int          `json:"role_id"`
	RoleName        string       `json:"role_name"`
	ModuleId        int          `json:"module_id"`
	ModuleName      string       `json:"module_name"`
	ModuleUrl       string       `json:"module_url"`
	ParentModuleId  int          `json:"parent_module_id"`
	ModuleOperation string       `json:"module_operation"`
	ModuleIcon      string       `json:"module_icon"`
	SonModuleFlag   bool         `json:"son_module_flag"`
	SonModule       []UserNavbar `json:"son_module"`
}

type PageOperat struct {
	IsAdd    int `json:"is_add"`
	IsUpdate int `json:"is_update"`
	IsDel    int `json:"is_del"`
	IsQuery  int `json:"is_query"`
}
