package models

type MemberInfo struct {
	Member_Id        int    `json:"member_id"`
	Member_Name      string `json:"member_name"`
	Member_Nick      string `json:"member_nick"`
	Member_Tel       string `json:"member_tel"`
	Member_WeChat_Id string `json:"member_wechat_id"`
	Member_Photo     string `json:"member_photo"`
	Create_Time      string `json:"create_time"`
}
