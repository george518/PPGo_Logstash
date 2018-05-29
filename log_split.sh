# /bin/bash

# 日志保存位置
base_path='/data/logs'
# 日志文件名称
log_name='access.log'
# ppgo_logstash 保存位置
ppgo_path='/root/soft/ppgo_logstash'


# 获取当前年信息和月信息
log_path=$(date -d yesterday +"%Y%m")
# 获取昨天的日信息
day=$(date -d yesterday +"%d")
# 按年月创建文件夹
mkdir -p $base_path/$log_path
# 备份昨天的日志到当月的文件夹
mv $base_path/$log_name $base_path/$log_path/$day.log

# 输出备份日志文件名
# echo $base_path/$log_path/$day.log
# 通过Nginx信号量控制重读日志 需要修改
kill -USR1 `cat /usr/local/nginx/logs/nginx.pid`
echo "\n" >> $base_path/$log_name

# 重启logprocess
cd $ppgo_path
./run.sh restart >> ./info.log