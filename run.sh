#!/bin/bash
# @Author: haodaquan
# @Date:   2018-05-29 17:44:45
# @Last Modified by:   haodaquan
# @Last Modified time: 2018-05-29 17:44:45


case $1 in
        start)
                nohup ./ppgo_logstash 2>&1 >> ./info.log 2>&1 /dev/null &
                echo "服务已启动..."
                sleep 1
        ;;
        stop)
                killall ppgo_logstash
                echo "服务已停止..."
                sleep 1
        ;;
        restart)
                killall ppgo_logstash
                sleep 1
                nohup ./ppgo_logstash 2>&1 >> ./info.log 2>&1 /dev/null &
                echo "服务已重启..."
                sleep 1
        ;;
        *)
                echo "$0 {start|stop|restart}"
                exit 4
        ;;
esac