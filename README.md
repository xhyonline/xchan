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

## 六、注意事项

当你使用第三方上传时(七牛云),后台界面上传成功的进度条不代表真实的上传进度。

例如:当你上传大文件时,也许你看到进度条已经100%了但是却迟迟没等到上传成功的回调。这是因为后台的进度条只代表前端将数据发送给了 xchan 服务的进度。但是如果使用七牛云,我们的服务端还要将数据二次转发给七牛存储。当转发成功后才会,才会有回调。(如下图所示:) 

![avatar](https://qiniu.xhyonline.com/ea89db385d02b9338669c37f3fde8897.png)

如果此时的你着急刷新了界面,其实这个上传过程并不会中断,也许你会觉得很卡,这是因为后端服务器还在上传上一条任务,因此占用了很大的IO,当然此时并不影响你继续继续上传。

如果你使用本地存储的模式,这个进度条就近乎与服务端存储实时了......

.......

## 七、上传超时问题解决

当你使用七牛云上传时,可能会经常遇到上传失败的可能。

我们存储的形式是是:   自定义域名 --> Nginx --> xchan 图床服务 --> 七牛云存储

此时可能就会出现几个问题:

1. 前端 JQuery 主动断开连接,这种表现形式为文件较大,JQuery 默认上传 30s 如果文件没有上传完毕,则会直接断开,最直观的现象就是下面这种

   ![avatar](https://qiniu.xhyonline.com/719399578b5fe7aeb6d00e61e12d7e5e.jpg)

解决方案:请直接修改代码,将前端代码中的 timeout 调高

![avatar](https://qiniu.xhyonline.com/3a66b8a62af3a792524ad3d1174de8bb.png)

2. 服务端超时

   有可能你会遇到我们的 xchan 应用程序与 nginx 之间的超时,因此你需要像下面这样的配置

![avatar](https://qiniu.xhyonline.com/8455c8ee946a79e9500ee8a918f5c8da.png)

```
server {
       listen 443 ssl;
        listen    [::]:443 ssl;
        server_name xchan.xhyonline.com;
        ssl_certificate /usr/local/ssl/5468484_xchan.xhyonline.com.pem;
        ssl_certificate_key /usr/local/ssl/5468484_xchan.xhyonline.com.key;
 	client_max_body_size     2048m;
	client_header_timeout    5m;
	client_body_timeout      5m;
        access_log /access.log;
        error_log  /var/logs/error.log  error;
	location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_connect_timeout      1500;   # 反向代理连接超时时长
	    proxy_send_timeout         1500;   # 代理发送超时时长
	    proxy_read_timeout         1500;   # 读取超时时长
	    proxy_http_version 1.1;         # 开启代理服务端长链接
	    proxy_set_header Connection "";  # 开启代理服务端长链接
            index  index.html index.htm;
        }
}

server {
    listen 80;
    server_name xchan.xhyonline.com;
    rewrite ^(.*)$ https://$host$1 permanent;
}

```





