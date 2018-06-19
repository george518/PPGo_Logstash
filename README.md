PPGo_Logstah
====
什么东西？What?
----
一款用于文件类日志实时监控的开源方案（说白了，就是读取文件日志，并将日志存储）。
读取文件日志---->存入InfluxDb---->Grafana查看数据

示例为
nginx日志:
监控指标：nginx流量，响应时间，访问地址等

基于golang


有什么价值？
----
1、监控nginx日志
2、influxDb连接池实现

用到了哪些？
----
1、golang1.8+
2、influxdb1.2
3、grafana4.2


效果展示
----
demo界面<br/>
![github](https://github.com/george518/PPGo_Logstash/blob/master/mock_data/images/demo.png?raw=true "github")
<br/><br/>


安装方法
----

linux

进入 https://github.com/george518/PPGo_Job/releases
下载 ppgo_job-linux-1.2.1.zip 并解压
进入文件夹，设置好数据库(创建数据库，导入ppgo_job.sql)和配置文件(conf/app.conf)
运行 ./run.sh start|stop

mac

进入https://github.com/george518/PPGo_Job/releases
下载 ppgo_job-mac-1.2.1.zip 并解压
进入文件夹，设置好数据库(创建数据库，导入ppgo_job.sql)和配置文件(conf/app.conf)
运行 ./run.sh start|stop


配置文件：

````
#GOLBAL
AppMode         = development  #development production
WebPort         = 8089
ReadNum         = 1 #Read进程数
ProcessNum      = 5
WriteNum        = 50
StorageType     = InfluxDb
LogType         = NginxLog

[InfluxDb]
Url         = http://118.89.238.78
Port        = 8086
User        = george
Pwd         = 123456
Name        = georgetest
Precision   = s
InitialCap  = 1 #连接池初始链接数量
MaxCap      = 3 #连接池最大连接数

[NginxLog]
Type        = nginx_log #不要修改
Table       = nginx_log
Path        = ./mock_data/access.log,./mock_data/access2.log #多个使用英文逗号隔开
Delimiter   = `||`
TimeFormat  = 02/Jan/2006:15:04:05 +0800
TimeLoc     = Asia/Shanghai
Tags        = `{"Path":{"id":"4","func":"split_str,url,trim_right_num","c_id":"1"},"Method":{"id":"4","func":"split_str,trim_left_1","c_id":"0"},"Schema":{"id":"4","func":"split_str,trim_right_1","c_id":"-1"},"Status":{"id":"5","func":""},"Ip":{"id":"0","func":""}}`
Fields      = `{"UpstreamTime":{"id":"10","func":"float64"},"RequestTime":{"id":"11","func":"float64"},"BytesSent":{"id":"6","func":"int"}}`
Times       = `{"id":"3","func":"trim_left_1,trim_right_1"}`

```
influxDb安装


log format
```
log_format  access  '$remote_addr||-||$remote_user||[$time_local]||"$request"||'
                  '$status||$body_bytes_sent||"$http_referer"||'
                  '"$http_user_agent"||"$http_x_forwarded_for"||'
                  '$upstream_response_time||$request_time';
```

联系我
----
qq群号:547564773
欢迎交流，欢迎提交代码。