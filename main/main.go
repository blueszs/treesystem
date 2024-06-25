package main

import (
	"fmt"
	"treesystem/common"
	"treesystem/common/cache"
	"treesystem/common/httpclient"
	"treesystem/models"
	"treesystem/sql"
)

var Conf models.Config
var CacheDb *cache.RedisCacheDb

func main() {
	/*ntp := func() {
		for {
			cf := common.DefaultNTPServerConf()
			t, err := common.Npt_Time(cf)
			if err != nil {
				fmt.Printf("获取NTP时间失败err:%s\n", err.Error())
			} else {
				fmt.Printf("%v\n", t)
				err = common.SetSystemDate(t)
				if err != nil {
					fmt.Printf("设置系统时间失败err:%s\n", err.Error())
				}
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}
	go ntp()*/
	sjson, err := common.ReadConfig("../config.json")
	if err != nil {
		fmt.Printf("config.json读取err：%s\n", err.Error())
		return
	}
	sjson = common.ReplaceByRegex("[\\s]", sjson, "")
	Conf, err = common.JsonDeserialize[models.Config](sjson)
	if err != nil {
		fmt.Printf("config.json序列化err：%s\n", err.Error())
		return
	}
	mysqlDB := sql.Init(Conf.Database)
	if len(Conf.Cache) > 0 {
		CacheDb = cache.InitRedisCache(Conf.Cache[0].IpAddr+":"+Conf.Cache[0].IpPort, Conf.Cache[0].Password, Conf.Cache[0].DbIndex)
	}
	err = mysqlDB.OpenDb()
	if err != nil {
		fmt.Printf("数据库打开err：#{err}\n")
		return
	}
	cookisM := make(map[string]string)
	/*cookisM["bbs_sid"] = common.GetRandomStr(26)
	cookisM["bbs_token"] = common.GetRandomStr(26)
	cookisM["cookie_test"] = common.GetRandomStr(26)*/
	body, _, err := httpclient.POST("https://gd.m.zjtcn.com/fac/getFacPrice.json", "application/json", cookisM, "{'code':'0051908900020'}")
	if err != nil {
		fmt.Printf("http request post err:#{err}\n")
	}
	bodystr := string(body)
	fmt.Println(bodystr)
	//colweb, _ := httpclient.GetCollect("https://v.douyin.com/i87P4yy/")
	//playsrc := help.GetRegexValue("playAddr%22%3A%5B%7B%22src%22%3A%22%2F%2F((?!(%7D))[\\S])+(%7D)", colweb)
	//fmt.Println(playsrc)
	//playsrc, _ = httpclient.URLEncode(playsrc)
	//fmt.Println(playsrc)
	//playsrc = help.GetRegexValue("[/]{2}[^\"]+", playsrc)
	//if len(playsrc) > 0 {
	//	playsrc = "https:" + playsrc
	//}
	//fmt.Println(playsrc)
	//colweb, _ = httpclient.GetCollect("https://www.google.com")
	//fmt.Println(colweb)
	//data, _ := dal.GetSystemConfigs(mysqlDB)
	//print(data)
	RunWebService(mysqlDB)
	/*
		colweb, _ := httpclient.GetCollect("https://v.douyin.com/UDJYsWU/")
		playsrc := help.GetRegexValue("playAddr%22%3A%5B%7B%22src%22%3A%22%2F%2F((?!(%7D))[\\S])+(%7D)", colweb)
		fmt.Println(playsrc)
		playsrc, _ = httpclient.URLEncode(playsrc)
		fmt.Println(playsrc)
		playsrc = help.GetRegexValue("[/]{2}[^\"]+", playsrc)
		if len(playsrc) > 0 {
			playsrc = "https:" + playsrc
		}
		fmt.Println(playsrc)
		liveweb, _ := httpclient.GetCollect("https://live.douyin.com/629463870615")
		livesrc := help.GetRegexValue("<script id=\"RENDER_DATA\" type=\"application/json\">(?!(</script>))[\\S]+</script>", liveweb)
		livesrc, _ = httpclient.URLEncode(livesrc)
		fmt.Println(livesrc)
		livesrc = help.GetRegexValue("(?<=(\"web_stream_url\":\\{\"flv_pull_url\":\\{\"FULL_HD1\":\"))http", livesrc)
		fmt.Println(livesrc)
		cookisM := make(map[string]string)
		//cookisM["bbs_sid"] = help.GetRandomStr(26)
		//cookisM["bbs_token"] = help.GetRandomStr(26)
		//cookisM["cookie_test"] = help.GetRandomStr(26)
		//dybody, _, _ := httpclient.GET("https://www.douyin.com/video/7228070658821950752", cookisM)
		//bystr := string(dybody)
		formsM := make(map[string]string)
		formsM["code"] = "0051908900020"
		body, _, err := httpclient.POST("https://gd.m.zjtcn.com/fac/getFacPrice.json", "application/json", cookisM, "{'code':'0051908900020'}")
		if err != nil {
			fmt.Printf("http request post err:#{err}\n")
		}
		bodystr := string(body)
		fmt.Println(bodystr)
		help.Trimming("C:\\Users\\Blue\\Pictures\\IMG_2747.jpg", "C:\\Users\\Blue\\dwhelper\\", 3, 3)
		err = encrypt.GenerateRSAKey(2048, "d:\\test\\", "2222")
		if err != nil {
			fmt.Printf("生成RSA密钥err:#{err}\n")
		}
		//var purl string
		//flag.StringVar(&purl, "url", "", "抓取url")
		sd, err := help.CreateQRBase64("https://ss0.baidu.com/94o3dSag_xI4khGko9WTAnF6hhy/baike/pic/item/c8177f3e6709c93df01d0b80913df8dcd00054dd.jpg")
		if err != nil {
			fmt.Printf("Create QRCode error:#{err}\n")
		}
		fmt.Println(sd)
		//fmt.Print("请输入获取网址：")
		//fmt.Scan(&purl)
		//flag.Parse()
		redisCl := cache.InitRedisCache("", "", 0)
		redisCl.Expire("K:2022-04-28:2", -1)
		mysqlDB.ExecSqlByTransact("insert into user_info(user_name,pass_word,nike_name,phone)values(?,?,?,?)", "root", "test123", help.ConvertCharToUnicode("测试"), "10086")
		query, _ := mysqlDB.QuerySql("select * from user_info limit 0,10")
		mysqlDB.CloseSqlConn()
		print(query)
		var cacheKey string
		for {
			cacheKey = "K:" + time.Now().Format("2006-01-02") + ":" + help.GetRandomStr(6)
			exFlag, _ := redisCl.Exists(cacheKey)
			if !exFlag {
				break
			}
		}
		mysqlDB.CloseSqlConn()
		//cookisM := make(map[string]string)
		//cookisM["bbs_sid"] = help.GetRandomStr(26)
		//cookisM["bbs_token"] = help.GetRandomStr(26)
		//cookisM["cookie_test"] = help.GetRandomStr(26)
		//body, resCookies, err := httpclient.GET("https://www.hifini.com", cookisM)
		//body, _, err := httpclient.GET(purl, cookisM)
		//body, _, err := httpclient.GET("https://www.hifini.com", cookisM)
		if err != nil {
			fmt.Println(err)
		} else {
			bodyImg, _, err2 := httpclient.GET("https://www.hifini.com/plugin/jan/img/bg-1.jpg", nil)
			if err2 != nil {
				fmt.Println(err2)
			}
			//help.Save(bodyImg, "d:\\temp\\test\\2022\\05\\05\\a.jpg", true)
			savePath := strings.Join([]string{"d:\\temp\\test", time.Now().Format("2006-01-02"), strings.Join([]string{help.GetRandomStr(8), "jpg"}, ".")}, "\\")
			b, err := help.BuffPhotoWatermarkText(bodyImg, "测试测试", ".\\resources\\font\\msyh_boot.ttf", savePath, 14, 100, 100)
			if err != nil {
				fmt.Printf("图片添加水印err：#{err}\n")
			}
			if !b {
				fmt.Print("图片添加水印失败\n")
			}
			//fmt.Println(string(body))
			shtml := string(body)
			shtml = help.FullToUnicode(shtml)
			err = redisCl.Set(cacheKey, shtml)
			if err != nil {
				fmt.Printf("设置Redis缓存键#{cacheKey}发生err,err:#{err}\n")
			}
			tmp := help.ConvertUnicodeToChar(shtml)
			fmt.Println(tmp)
			fmt.Println("hello world")
			// 关闭redis客户端链接
			defer redisCl.Close()
		}*/

	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, tpl.(1, user))
	})
	fmt.Println("服务端口:8088")                 //控制台输出信息
	err := http.ListenAndServe(":8088", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}*/
}
