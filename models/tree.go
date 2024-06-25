package models

type TreeDetail struct {
	Detail_Id         int               `json:"detail_id"`
	Tree_Id           int               `json:"tree_id"`
	Tree_Name         string            `json:"tree_name"`
	Plant_Id          int               `json:"plant_id"`
	Plant_Name        string            `json:"plant_name"`
	Detail_Dim_Key1   string            `json:"detail_dim_key1"`
	Detail_Dim_Value1 string            `json:"detail_dim_value1"`
	Detail_Dim_Key2   string            `json:"detail_dim_key2"`
	Detail_Dim_Value2 string            `json:"detail_dim_value2"`
	Detail_Dim_Key3   string            `json:"detail_dim_key3"`
	Detail_Dim_Value3 string            `json:"detail_dim_value3"`
	Detail_Dim_Key4   string            `json:"detail_dim_key4"`
	Detail_Dim_Value4 string            `json:"detail_dim_value4"`
	Detail_Dim_Key5   string            `json:"detail_dim_key5"`
	Detail_Dim_Value5 string            `json:"detail_dim_value5"`
	Active_Price      float64           `json:"active_price"`
	Real_Price        float64           `json:"real_price"`
	Prime_Price       float64           `json:"prime_price"`
	Store_Quantity    int               `json:"store_quantity"`
	Sales_Quantity    int               `json:"sales_quantity"`
	Tree_Unity        string            `json:"tree_unity"`
	Flag              int               `json:"flag"`
	Flag_Info         string            `json:"flag_info"`
	Insert_Time       string            `json:"insert_time"`
	TreeDetailPhotos  []TreeDetailPhoto `json:"TreeDetailPhotos"`
	Tree_Attr         []TreeAttr        `json:"tree_attr"`
}

type TreeAttr struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TreeDetailPhoto struct {
	Id            int    `json:"id"`
	Detail_Id     int    `json:"detail_id"`
	Img_Source    string `json:"img_source"`
	Img_Thumbnail string `json:"img_thumbnail"`
	Sort_Num      int    `json:"sort_num"`
}

type TreePlant struct {
	Plant_Id       int    `json:"plant_id"`
	Plant_Name     string `json:"plant_name"`
	Plant_Address  string `json:"plant_address"`
	Plant_Tel      string `json:"plant_tel"`
	Plant_Info     string `json:"plant_info"`
	Plant_Contacts string `json:"plant_contacts"`
	Flag           int    `json:"flag"`
	Flag_Info      string `json:"flag_info"`
	Insert_Time    string `json:"insert_time"`
}

type TreeClass struct {
	Class_Id          int    `json:"class_id"`
	Class_Name        string `json:"class_name"`
	Parent_Class_Id   int    `json:"parent_class_id"`
	Parent_Class_Name string `json:"parent_class_name"`
	Insert_Time       string `json:"insert_time"`
}

type TreeInfo struct {
	Tree_Id                   int             `json:"tree_id"`
	Tree_Name                 string          `json:"tree_name"`
	Class_Id                  int             `json:"class_id"`
	Class_Name                string          `json:"class_name"`
	Tree_Info                 string          `json:"tree_info"`
	Short_Tree_Info           string          `json:"short_tree_info"`
	Tree_Main_Photo           string          `json:"tree_main_photo"`
	Tree_Main_Photo_Thumbnail string          `json:"tree_main_photo_thumbnail"`
	Flag                      int             `json:"flag"`
	Flag_Info                 string          `json:"flag_info"`
	Insert_Time               string          `json:"insert_time"`
	Tree_Photos               []string        `json:"tree_photos"`
	Details                   []TreeDetail    `json:"details"`
	Attribute                 []TreeAttribute `json:"attribute"`
}

type TreeAttribute struct {
	Name   string     `json:"name"`
	Values []AtrrInfo `json:"values"`
}

type AtrrInfo struct {
	Attrdetal string `json:"attrdetal"`
	Attrsate  string `json:"attrsate"`
}
