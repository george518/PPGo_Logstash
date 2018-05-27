/************************************************************
** @Description: PPGo_LogProcess
** @Author: george hao
** @Date:   2018-05-14 17:03
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:03
*************************************************************/
package main

import (
	"log"
	"time"

	"github.com/george518/PPGo_Logstash/config"
	"github.com/george518/PPGo_Logstash/logdig"
	"github.com/george518/PPGo_Logstash/monitor"
	"github.com/george518/PPGo_Logstash/process"
	. "github.com/george518/PPGo_Logstash/storage"
	"github.com/george518/PPGo_Logstash/types"
)

func main() {

	var poolConfig = &PoolConfig{
		InitialCap: 1,
		MaxCap:     3,
		Factory:    DbFactory,
		Close:      DbClose,
		//链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
		IdleTimeout: 5 * time.Second,
	}

	dbpool, err := NewChannelPool(poolConfig)
	if err != nil {
		log.Fatal("create db pool error")
	}

	//for i := 0; i < 10; i++ {
	//	conn, err := dbpool.Get()
	//	if err != nil {
	//		log.Fatal("get conn error")
	//	}
	//	fmt.Println(dbpool.Len())
	//	fmt.Println(conn)
	//	dbpool.Put(*conn)
	//	time.Sleep(1 * time.Second)
	//}

	wc := make(chan *types.Message, 200)
	rc := make(chan []byte, 200)
	Conf := config.LoadConfig()

	logData := &logdig.LogData{
		Rc:   rc,
		Path: Conf.LogInfo.Path,
	}

	logProcess := &process.LogProcess{
		Wc:      wc,
		Rc:      rc,
		LogInfo: Conf.LogInfo,
	}

	storage := &Storage{
		Wc:  wc,
		Db:  Conf.StorageDb,
		Env: Conf.AppMode,
	}

	readNum := Conf.ReadNum
	procNum := Conf.ProcessNum
	writeNum := Conf.WriteNum

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
		WebPort:   Conf.WebPort,
	}
	m.Start(logProcess)

}
