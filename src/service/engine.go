package service

import (
	"ant_god/src/domain"
	"context"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

//将一串cookie分类整合成数组
func handleCookies(cookies string) map[string]string {
	result := make(map[string]string)
	cookie := strings.Split(cookies, ";")
	for _, value := range cookie {
		item := strings.Split(value, "=")
		if len(item) == 2 {
			result[strings.TrimSpace(item[0])] = strings.TrimSpace(item[1])
		}
	}
	return result
}

//适配需要验证的代理服务器
func authProxyServer(username, password string, taskCtx context.Context) {
	//监听中断事件
	lctx, lcancel := context.WithCancel(taskCtx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch ev := ev.(type) {
		//若由于请求停止引起的中断，则继续请求
		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(taskCtx, fetch.ContinueRequest(ev.RequestID))
			}()
		//若由于代理服务器需要验证而引起的中断，则验证。
		case *fetch.EventAuthRequired:
			if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				go func() {
					err := chromedp.Run(taskCtx,
						fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
							Response: fetch.AuthChallengeResponseResponseProvideCredentials,
							Username: username,
							Password: password,
						}),
						fetch.Disable(),
					)
					if err != nil {
						log.Fatalln(err)
					}
					lcancel()
				}()
			}
		}
	})
}

//初始化
func InitPatch(configuration domain.Configuration) (context.Context, []context.CancelFunc) {
	cancel := make([]context.CancelFunc, 3)

	chromeExecPath := configuration.ChromeExecPath
	proxyServer := ""
	userAgent := ""
	headless := true
	timeout := time.Second * 2000
	cookies := ""
	if configuration.Cookies != "" {
		cookies = configuration.Cookies
	}
	if configuration.ProxyServer.Address != "" {
		proxyServer = configuration.ProxyServer.Address
	}
	if configuration.UserAgent != "" {
		userAgent = configuration.UserAgent
	}
	if configuration.Timeout != 0 {
		timeout = time.Second * time.Duration(configuration.Timeout)
	}
	if !configuration.Headless {
		headless = false
	}

	//指定分配器选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromeExecPath),
		chromedp.ProxyServer(proxyServer),
		chromedp.UserAgent(userAgent),
		chromedp.Flag("headless", headless),
		chromedp.Flag("blink-settings", "imagesEnable=false"),
	)

	//新建分配器作为Context
	allocatorCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	cancel[0] = cancel1

	// 给每个页面的爬取设置超时时间
	timeoutCtx, cancel2 := context.WithTimeout(allocatorCtx, timeout)
	cancel[1] = cancel2

	//设置日志打印
	taskCtx, cancel3 := chromedp.NewContext(timeoutCtx, chromedp.WithLogf(log.Printf))
	cancel[2] = cancel3

	//监听代理服务器的验证请求，并验证，若不需要验证，传空即可。
	authProxyServer(configuration.ProxyServer.Username, configuration.ProxyServer.Password, taskCtx)

	//第一次调用Run方法会创建chrome进程，Run方法的可变参数actions是有顺序的。
	if err := chromedp.Run(taskCtx,
		//ActionFunc用来自定义操作，比如可以用来设置Cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			result := handleCookies(cookies)
			if len(result) != 0 {
				for key, value := range result {
					err := network.SetCookie(key, value).Do(ctx)
					if err != nil {
						log.Fatalln(err)
					}
				}
				return nil
			} else {
				return nil
			}
		}), page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
	); err != nil {
		log.Fatalln(err)
	}
	return taskCtx, cancel
}

// Patch 爬取单页面
func Patch(url string, taskCtx context.Context) string {
	var html string

	//开始爬取数据
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(`document.querySelector("html")`, &html, chromedp.ByJSPath),
	); err != nil {
		//浏览器自身的错误，则跳过
		if strings.Contains(err.Error(), "page load error") {
			return "给定页面无法的打开"
		} else {
			log.Fatalln(err)
		}
	}
	return html
}

// PatchAll 爬取多页面
func PatchAll(urls []string, taskCtx context.Context) []string {
	htmls := make([]string, len(urls))
	for index, value := range urls {
		html := Patch(value, taskCtx)
		htmls[index] = html
	}
	return htmls
}

func DestroyContext(cancle []context.CancelFunc) {
	for _, value := range cancle {
		value()
	}
}
