package main

import (
	"ant_god/src/dao"
	"ant_god/src/ext/xray"
	"ant_god/src/service"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

//互斥锁
var mutex sync.Mutex

//等待组
var wg sync.WaitGroup

//工人
func Worker(id int, jobs <-chan int, allURLs []string, taskContext context.Context, domain string) {
	for job := range jobs {
		//开始打工
		fmt.Println("worker", id, "started  job", job)

		//操作共享数据前加锁
		mutex.Lock()
		tempHtml := service.Patch(allURLs[job], taskContext)
		//操作完解锁
		mutex.Unlock()

		//数据处理
		tempURLs := service.DataProcess(tempHtml, domain)

		if tempURLs != nil {
			//操作共享数据前加锁
			mutex.Lock()
			for _, value2 := range tempURLs {
				allURLs = append(allURLs, value2)
			}
			//操作完解锁
			mutex.Unlock()

			//打工结束
			fmt.Println("worker", id, "finished job", job)
			wg.Done()
		} else {
			//打工结束
			fmt.Println("worker", id, "finished job", job)
			wg.Done()
		}
	}

}

func main() {

	//1.从命令行读取URL
	url, filePath := dao.ReadCmdOptions()

	if *url == "" {
		log.Fatalln("-u选项为必选")
	}
	if strings.Contains(*url, "@") {
		log.Fatalln("url中不能包含@符号")
	}
	if !strings.Contains(*url, "https") && !strings.Contains(*url, "http") {
		log.Fatalln("url中必须包含协议")
	}

	//2.读取配置文件内容
	configuration := dao.ReadConfiguration()

	//3.初始化爬虫引擎
	taskContext, cancel := service.InitPatch(configuration)

	var allURLs []string
	domain := service.ParseURL(*url)

	//4.循环爬取
	//4.1 获取页面
	html := service.Patch(*url, taskContext)

	//4.2 从页面中提取URL
	urls := service.DataProcess(html, domain)

	if urls != nil {
		//4.3 将提取的url汇总
		for _, value := range urls {
			allURLs = append(allURLs, value)
		}
		//4.4 数据去重
		allURLs = service.RemoveRepeatedElement(allURLs)
		//4.5 递归爬取子页面
		start := 0
		end := len(allURLs)
		jobs := make(chan int, end)

		for i := 0; i < configuration.Worker; i++ {
			//三个打工仔
			go Worker(i, jobs, allURLs, taskContext, domain)
		}

		for {
			for i := start; i < end; i++ {
				wg.Add(1)
				jobs <- i
			}
			wg.Wait()
			//数据去重
			allURLs = service.RemoveRepeatedElement(allURLs)
			if len(allURLs) != end {
				start = end
				end = len(allURLs)
				jobs = make(chan int, end)
			} else {
				break
			}
		}

		//将所有结果保存到文件
		dao.SaveFile(allURLs, *filePath)

	}
	//5.是否启用xray扫描
	if configuration.XrayExt.Enabled {
		fmt.Println(configuration.XrayExt.XrayPath)
		if configuration.XrayExt.XrayPath == "" {
			log.Fatalln("-----------------------未指定xray路径-----------------------")
		} else {
			xray.StartScan(configuration.XrayExt.XrayPath, *filePath)
		}
	}
	service.DestroyContext(cancel)
}
