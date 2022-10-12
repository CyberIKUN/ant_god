## 介绍

🐶这是一款爬虫器，有以下特点：

- 自动爬取js渲染后的url页面里的url链接
- 链接去重
- 链接限制同域
- 递归爬取子页面
- 将数据交给xray进行漏扫

这只是一个学习爬虫的案例，代码里带有了注释，供大家对js渲染的数据爬取进行学习。

## 使用

1. 😪`git clone https://github.com/CyberIKUN/ant_god.git`
2. 😏到项目路径执行`./ant_god.exe -u https://www.baidu.com/`

## 配置文件

在src\conf目录下有个conf.yaml配置文件，可根据配置文件进行参数选调。

## 命令行参数

只有两个：

1. -u：指定url
2. -o：指定输出文件路径，例如：-u C:\Users\1.txt