package models

type ListenClass struct {
	Class_Id   int    `json:"class_id"`
	Class_Name string `json:"class_name"`
	Flag       int    `json:"flag"`
	Flag_Info  string `json:"flag_info"`
}

type ListenSeries struct {
	Series_Id   int    `json:"series_id"`
	Series_Name string `json:"series_name"`
	Class_Id    int    `json:"class_id"`
	Class_Name  string `json:"class_name"`
	Series_Info string `json:"series_info"`
	Series_Img  string `json:"series_img"`
	Add_Time    string `json:"add_time"`
	Flag        int    `json:"flag"`
	Flag_Info   string `json:"flag_info"`
}
type ListenInfo struct {
	Listen_Id   int    `json:"listed_id"`
	Listen_Name string `json:"listed_name"`
	Series_Id   int    `json:"series_id"`
	Series_Name string `json:"series_name"`
	Listen_Info string `json:"listed_info"`
	Listen_Url  string `json:"listed_url"`
	Listen_Img  string `json:"listed_img"`
	Listen_Ep   int    `json:"listed_ep"`
	Add_Time    string `json:"add_time"`
	Flag        int    `json:"flag"`
	Flag_Info   string `json:"flag_info"`
}
