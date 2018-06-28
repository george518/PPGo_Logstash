/************************************************************
** @Description: logstash
** @Author: george hao
** @Date:   2018-06-15 15:25
** @Last Modified by:  george hao
** @Last Modified time: 2018-06-15 15:25
*************************************************************/
package logstash

import (
	"github.com/george518/PPGo_Logstash/config"
	"github.com/george518/PPGo_Logstash/logdig"
	"github.com/george518/PPGo_Logstash/monitor"
	"github.com/george518/PPGo_Logstash/process"
	. "github.com/george518/PPGo_Logstash/storage"
	"github.com/george518/PPGo_Logstash/types"
	"log"
	"os"
	"time"
)

func Run(configFile string, exitChan chan int) {
	Conf := config.LoadConfig(configFile)

	wc := make(chan *types.LogMessage, 200)
	rc := make(chan []byte, 200)

	//monitor file
	//监控日志记录
	dir, _ := os.Getwd() //当前的目录
	types.MonitorLogPath = dir + "/monitor.log"
	_, err := os.Create(types.MonitorLogPath)
	if err != nil {
		log.Println(" monitor log file create err", err)
		return
	}
	//db pool
	InitialCap, err := Conf.Storage.Key("InitialCap").Int()
	if err != nil {
		InitialCap = 1
	}
	MaxCap, err := Conf.Storage.Key("MaxCap").Int()
	if err != nil {
		MaxCap = 3
	}
	var poolConfig = &PoolConfig{
		InitialCap: InitialCap,
		MaxCap:     MaxCap,
		Factory:    DbFactory,
		Close:      DbClose,
		//链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
		IdleTimeout: 5 * time.Second,
		//数据库链接
		Conf: Conf.Storage,
	}

	dbpool, err := NewChannelPool(poolConfig)
	if err != nil {
		log.Fatal("create db pool error")
	}

	logData := &logdig.LogData{
		Rc:   rc,
		Conf: Conf.LogType,
	}

	logProcess := &process.LogProcess{
		Wc:      wc,
		Rc:      rc,
		LogInfo: Conf.LogType,
	}

	writeNum, err := Conf.Storage.Key("WriteNum").Int()
	if err != nil {
		log.Println(" WriteNum error")
		writeNum = 1
	}
	storage := &Storage{
		Wc:       wc,
		Table:    Conf.LogType.Key("Table").String(),
		Env:      Conf.Global.Key("AppMode").String(),
		WriteNum: writeNum,
	}

	logData.Read()
	logProcess.Process()
	storage.Save(*dbpool)

	m := &monitor.Monitor{
		StartTime: time.Now(),
		Data:      types.SystemInfo{},
		WebPort:   Conf.Global.Key("WebPort").String(),
	}
	m.Start(logProcess)

	select {
	case n := <-types.ExitChan:
		log.Println("ppgo_logstash is stoped:", n)

		m.PrintMonitor(logProcess)
		os.Exit(0)
	}

}
