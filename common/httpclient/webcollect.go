package httpclient

import (
	"context"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

// GetCollect http请求
// @curl http请求的地址
func GetCollect(curl string) (string, error) {
	timeCtx, cancel := GetChromeCtx(false, curl)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(timeCtx,
		// 导航到目标URL
		chromedp.Navigate(curl),
		//chromedp.WaitReady(`document`),
		// 等待某个特定的元素加载完成，确保页面JavaScript执行完毕
		chromedp.WaitVisible(`html`, chromedp.ByQuery),
		// 获取页面的HTML源码
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("读取页面失败1：", err.Error())
		return "", err
	}
	return htmlContent, nil
}

// URLEncode Url解码
// @content 待解码字符串
func URLEncode(content string) (string, error) {
	htmlContent, err := url.PathUnescape(content)
	if err != nil {
		log.Println("读取页面失败1：", err.Error())
		return "", err
	}
	return htmlContent, nil
}

// 检查是否有9222端口，来判断是否运行在linux上
func checkChromePort() bool {
	addr := net.JoinHostPort("", "9222")
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// ChromeCtx 使用一个实例
var ChromeCtx context.Context

func GetChromeCtx(focus bool, curl string) (context.Context, context.CancelFunc) {
	if ChromeCtx == nil || focus {
		allocOpts := chromedp.DefaultExecAllocatorOptions[:]
		allocOpts = append(allocOpts,
			chromedp.DisableGPU,
			chromedp.Flag("headless", true),                        // headless模式，“true”隐藏chrome，“false”显示chrome浏览器（debug使用）
			chromedp.Flag("blink-settings", "imagesEnabled=false"), // 禁止图片加载（可减少流量使用和响应时间）
			chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36`),
			chromedp.Flag("accept-language", `zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6`),
		)
		if checkChromePort() {
			ChromeCtx, _ = chromedp.NewRemoteAllocator(context.Background(), curl)
		} else {
			ChromeCtx, _ = chromedp.NewExecAllocator(context.Background(), allocOpts...)
		}
	}
	ctx, _ := chromedp.NewContext(ChromeCtx)
	// 给每个页面的爬取设置超时时间
	timeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	return timeCtx, cancel
}
