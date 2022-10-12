package dao

import (
	"ant_god/src/domain"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// ReadConfiguration 读取配置文件
func ReadConfiguration() domain.Configuration {
	configuration := domain.Configuration{}

	//读文件
	file, err := ioutil.ReadFile(`src/conf/conf.yaml`)
	if err != nil {
		log.Fatalln(err)
	}

	//反序列化
	err = yaml.Unmarshal(file, &configuration)
	if err != nil {
		log.Fatalln(err)
	}

	//指定默认值
	if configuration.ChromeExecPath == "" {
		configuration.ChromeExecPath = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	}

	if configuration.Timeout == 0 {
		configuration.Timeout = 30
	}
	if configuration.Worker == 0 {
		configuration.Worker = 3
	}
	return configuration
}
