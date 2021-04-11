# xchan
xchan 是一个用Golang开发的对象存储工具,适合供个人开发者使用

支持以下两种上传模式

1. 本地存储
2. 七牛云存储
3. 生成 MarkDown 形式的地址

## 一、安装引导界面

![avatar](https://qiniu.xhyonline.com/4fe28c7a7dc41eee03ec465f7e073242)

![avatar](https://qiniu.xhyonline.com/7b0064f739345d84e8284a26be5c4b38)

## 二、拖拽上传

![avatar](https://qiniu.xhyonline.com/bc934e51ac944209a3d3fa04142b29bf)

## 三、文件管理

![avatar](https://qiniu.xhyonline.com/9bc300d1d02051dc28d26308a8b94708)

## 四、存储切换

![avatar](https://qiniu.xhyonline.com/c6e2b1ba4e497930e80a9b893e6b47d5)

![avatar](https://qiniu.xhyonline.com/8fde4966c7703229b30e8ba8fc0a8115)

## 五、安装方式

1. 请在 Linux 系统下下载该源码,并且安装 Golang 环境,设置 GoProxy 

2. 切到项目根目录执行命令:

   ```
   go build
   ```

   切记不要 使用 `go run main.go` ,您应该直接编译,因为此操作会导致存储的文件再 /tmp目录

3. 运行编译后的文件,项目会启动在 80 端口,您可以执行`./xchan -p 8089`切换端口,且切换后自行用 Nginx 反向代理的形式配置
4. 直接访问即可进入安装引导界面



