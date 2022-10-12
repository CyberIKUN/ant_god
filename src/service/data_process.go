package service

import (
	"regexp"
	"strings"
	"unicode"
)

/*
	将html页面中的url数据进行处理：
		1.提取url数据
		2.去除不同域
*/
func DataProcess(data string, domain string) []string {
	if data != "" {
		//提取规则
		re := regexp.MustCompile(`(src|href)="([._:/?%&=0-9a-zA-Z^"]+)"`)
		result := re.FindAllStringSubmatch(data, -1)
		if result != nil {
			var urls []string

			//url数据处理
			for _, value := range result {
				//去除只有//或/的url
				if value[2] == "//" || value[2] == "/" {
					continue
				}
				//当点开这些图片页面，浏览器会卡在渲染，导致爬虫停止
				if value[2][len(value[2])-3:] == "png" ||
					value[2][len(value[2])-3:] == "jpg" ||
					value[2][len(value[2])-3:] == "svg" ||
					value[2][len(value[2])-4:] == "jpeg" ||
					value[2][len(value[2])-3:] == "gif" {
					continue
				}
				//没有http或https前缀的处理：
				//判断是否同域，然后拼接
				if !strings.HasPrefix(value[2], "http") && !strings.HasPrefix(value[2], "https") {
					if strings.HasPrefix(value[2], "//") {
						patchDomain := ParseURL(value[2])
						if patchDomain == domain {
							urls = append(urls, "https:"+value[2])
						} else {
							continue
						}
					} else if strings.HasPrefix(value[2], "/") {
						urls = append(urls, "https://"+domain+value[2])
					} else if strings.HasPrefix(value[2], ":") {
						if unicode.IsDigit(rune(value[2][strings.Index(value[2], ":")+1])) {
							urls = append(urls, "https://"+domain+"/"+value[2])
						} else {
							continue
						}
					}
					//有http或https前缀的处理：
					//判断是否同域
				} else {
					patchDomain := ParseURL(value[2])
					if patchDomain == domain {
						urls = append(urls, value[2])
					} else {
						continue
					}
				}
			}
			return urls
		} else {
			return nil
		}
	} else {
		return nil
	}
}

// ParseURL 从URL中提取域名
func ParseURL(url string) string {
	var domain string
	//"https://www.baidu.com"
	if strings.Contains(url, "http") || strings.Contains(url, "https") {
		frontSep := url[strings.Index(url, ":")+3:]
		if !strings.Contains(frontSep, "/") {
			domain = frontSep
		} else {
			domain = frontSep[:strings.Index(frontSep, "/")]
		}
		//"//www.baidu.com/
	} else {
		frontSep := url[strings.Index(url, "/")+2:]
		if !strings.Contains(frontSep, "/") {
			domain = frontSep
		} else {
			domain = frontSep[:strings.Index(frontSep, "/")]
		}
	}

	return domain
}

// RemoveRepeatedElement 数据去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
