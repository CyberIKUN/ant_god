package xray

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// StartScan xray扫描模块
func StartScan(xrayPath string, filePath string) {
	urls := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	year, month, day := time.Now().Date()
	fileName := strconv.Itoa(year) + "_" + month.String() + "_" + strconv.Itoa(day) + ".html"
	for _, value := range urls {
		cmd := exec.Command(xrayPath, "webscan", "--url", value, "--html-output", "out/"+fileName)
		//错误输出使用标准输出
		cmd.Stderr = cmd.Stdout
		//拿到命令行输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalln(err)
		}
		//开始执行命令，不等待命令执行结果
		if err = cmd.Start(); err != nil {
			log.Fatalln(err)
		}
		//循环读取命令行输出
		for {
			tmp := make([]byte, 1024)
			//实时获取输出
			_, err := stdout.Read(tmp)
			//读不到输出时退出，即执行完毕
			if err != nil {
				break
			}
			fmt.Println(string(tmp))
		}
		//等待执行完成
		if err = cmd.Wait(); err != nil {
			log.Fatalln(err)
		}
	}
}
