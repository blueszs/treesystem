package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GET http的GET请求
// @curl http请求的地址
// @cookiesMap http请求的cookie字典对象
func GET(curl string, cookiesMap map[string]string) ([]byte, map[string]string, error) {
	client := &http.Client{}
	//建立一个请求
	req, err := http.NewRequest("GET", curl, nil)
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	//设置http请求头协议
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	// 添加Cookie
	for form := range cookiesMap {
		req.AddCookie(&http.Cookie{Name: form, Value: cookiesMap[form]})
	}
	resp, err := client.Do(req) //提交
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	defer ResponseClose(resp)
	cookies := resp.Cookies()
	//遍历cookies
	respCookies := make(map[string]string)
	for _, cookie := range cookies {
		respCookie := cookie.String()
		kv := strings.Split(respCookie, "=")
		if len(kv) > 1 {
			respCookies[kv[0]] = kv[1]
		} else {
			respCookies[kv[0]] = ""
		}
	}
	// 读取响应内容
	body, err1 := io.ReadAll(resp.Body)
	return body, respCookies, err1
}

// POSTForm http的POST请求
// @curl http请求的地址
// @cookiesMap http请求的cookie字典对象
// @formMap POST的Form参数
func POSTForm(curl string, cookiesMap map[string]string, formMap map[string]string) ([]byte, map[string]string, error) {
	// 表单数据
	formData := url.Values{}
	//拼接form请求参数
	for form := range formMap {
		formData.Add(form, formMap[form])
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", curl, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	//设置http请求头协议
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 添加Cookie
	for form := range cookiesMap {
		req.AddCookie(&http.Cookie{Name: form, Value: cookiesMap[form]})
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	defer ResponseClose(resp)
	cookies := resp.Cookies() //遍历cookies
	respCookies := make(map[string]string)
	for _, cookie := range cookies {
		respCookie := cookie.String()
		kv := strings.Split(respCookie, "=")
		if len(kv) > 1 {
			respCookies[kv[0]] = kv[1]
		} else {
			respCookies[kv[0]] = ""
		}
	}
	// 读取响应内容
	body, err1 := io.ReadAll(resp.Body)
	//fmt.Println(string(body)) //网页源码
	return body, respCookies, err1
}

// POST http的POST请求
// curl http请求的地址
// reqType http请求类型（"application/json"）
// cookiesMap http请求的cookie字典对象
// reqBody POST的收据（json，文件等）
func POST(curl, reqType string, cookiesMap map[string]string, reqBody any) ([]byte, map[string]string, error) {

	body, _ := json.Marshal(reqBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", curl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	//添加http请求头协议
	req.Header.Set("Content-Type", reqType)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;text/json;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	// 添加Cookie
	for form := range cookiesMap {
		req.AddCookie(&http.Cookie{Name: form, Value: cookiesMap[form]})
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("fatal error:%v\n", err)
		return nil, nil, err
	}
	defer ResponseClose(resp)
	cookies := resp.Cookies() //遍历cookies
	respCookies := make(map[string]string)
	for _, cookie := range cookies {
		respCookie := cookie.String()
		kv := strings.Split(respCookie, "=")
		if len(kv) > 1 {
			respCookies[kv[0]] = kv[1]
		} else {
			respCookies[kv[0]] = ""
		}
		fmt.Printf("cookie:%s\n", respCookie)
	}
	// 读取响应内容
	body, err1 := io.ReadAll(resp.Body)
	//fmt.Println(string(body)) //网页源码
	return body, respCookies, err1
}

// ResponseClose http response响应关闭
func ResponseClose(b *http.Response) {
	if b != nil {
		err := b.Body.Close()
		if err != nil {
			fmt.Printf("http response close failed:%v\n", err)
		} else {
			fmt.Println("http response close success!")
		}
	} else {
		fmt.Println("http response is null or undefined!")
	}
}
