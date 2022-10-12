package dao

import "flag"

// ReadCmdOptions 读取命令行选项
func ReadCmdOptions() (*string, *string) {
	//-u选项指示爬取的url，必选
	url := flag.String("u", "", "Set up crawling links")
	//-o选项指定输出文件路径，可选
	filePath := flag.String("o", "out/out.txt", "file retention path")
	flag.Parse()
	return url, filePath
}
