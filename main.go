/************************************************************
** @Description: PPGo_LogProcess
** @Author: george hao
** @Date:   2018-05-14 17:03
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:03
*************************************************************/
package main

import (
	"flag"
	"github.com/george518/PPGo_Logstash/logstash"
	"github.com/george518/PPGo_Logstash/types"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configFile *string = flag.String("c", "./config/conf.ini", "Use -c <config_file_path>")

func main() {

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2:
				log.Println(" Ready to quit ", s)
				types.ExitReadChan <- 1
			default:
				log.Println("other", s)
			}
		}
	}()

	log.Println(" ppgo_logstash is starting...")
	//开始愉快的干活啦
	logstash.Run(*configFile, types.ExitReadChan)

}
