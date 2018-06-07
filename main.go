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
	"github.com/george518/PPGo_Logstash/config"
	"github.com/george518/PPGo_Logstash/logdig"
	"github.com/george518/PPGo_Logstash/monitor"
	"github.com/george518/PPGo_Logstash/process"
	. "github.com/george518/PPGo_Logstash/storage"
	"github.com/george518/PPGo_Logstash/types"
	"log"
	"time"
)

var config_file *string = flag.String("c", "./config/conf.ini", "Use -c <config_file_path>")

func main() {
	flag.Parse()
	Conf := config.LoadConfig(*config_file)

	wc := make(chan *types.LogMessage, 200)
	rc := make(chan []byte, 200)

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

	storage := &Storage{
		Wc:    wc,
		Table: Conf.LogType.Key("Table").String(),
		Env:   Conf.Global.Key("AppMode").String(),
	}

	readNum, _ := Conf.Global.Key("ReadNum").Int()
	procNum, _ := Conf.Global.Key("ProcessNum").Int()
	writeNum, _ := Conf.Global.Key("WriteNum").Int()

	for i := 0; i < readNum; i++ {
		go logData.Read()
	}

	//log.Println(storage)
	for i := 0; i < procNum; i++ {
		go logProcess.Process()
	}

	for i := 0; i < writeNum; i++ {
		go storage.Save(*dbpool)
	}

	m := &monitor.Monitor{
		StartTime: time.Now(),
		Data:      types.SystemInfo{},
		WebPort:   Conf.Global.Key("WebPort").String(),
	}
	m.Start(logProcess)

}
