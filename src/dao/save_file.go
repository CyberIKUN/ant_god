package dao

import (
	"log"
	"os"
)

//保存文件
func SaveFile(urls []string, filePath string) {
	//打开文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	//写入值
	for _, value := range urls {
		file.WriteString(value + "\n")
	}
}
