/************************************************************
** @Description: PPGo_LogProcess
** @Author: george hao
** @Date:   2018-05-14 17:03
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:03
*************************************************************/
package main

import (
	"github.com/george518/PPGo_Logstash/config"
	"github.com/george518/PPGo_Logstash/logdig"
	"github.com/george518/PPGo_Logstash/monitor"
	"github.com/george518/PPGo_Logstash/process"
	"github.com/george518/PPGo_Logstash/storage"
	"github.com/george518/PPGo_Logstash/types"
	"github.com/go-ini/ini"
	"time"
)

type WrCh struct {
	Wc chan *process.Message
	Rc chan []byte
}

var Cfg ini.File

var Conf *config.ConfigGlobal

func init() {
	Cfg = config.Load()
}

func main() {

	wc := make(chan *process.Message, 200)
	rc := make(chan []byte, 200)

	wrch := WrCh{
		Wc: wc,
		Rc: rc,
	}
	logData := &logdig.LogData{
		Rc:   wrch.Rc,
		Path: Cfg.Section("log_info").Key("log_path").String(),
	}

	logProcess := &process.LogProcess{
		Wc:         wrch.Wc,
		Rc:         wrch.Rc,
		TimeLoc:    Cfg.Section("").Key("time_loc").String(),
		Regexp:     Cfg.Section("log_info").Key("log_regexp").String(),
		TimeFormat: Cfg.Section("log_info").Key("log_time").String(),
	}

	storage := &storage.Storage{
		Wc:          wrch.Wc,
		DbUrl:       Cfg.Section("storage").Key("db_url").String(),
		DbPort:      Cfg.Section("storage").Key("db_port").String(),
		DbUser:      Cfg.Section("storage").Key("db_user").String(),
		DbPwd:       Cfg.Section("storage").Key("db_pwd").String(),
		DbName:      Cfg.Section("storage").Key("db_name").String(),
		DbPrecision: Cfg.Section("storage").Key("db_precision").String(),
		DbTable:     Cfg.Section("storage").Key("db_table").String(),
		Env:         Cfg.Section("").Key("app_mode").String(),
	}

	readNum, _ := Cfg.Section("").Key("read_num").Int()
	procNum, _ := Cfg.Section("").Key("process_num").Int()
	writeNum, _ := Cfg.Section("").Key("write_num").Int()

	for i := 0; i < readNum; i++ {
		go logData.Read()
	}

	//log.Println(storage)
	for i := 0; i < procNum; i++ {
		go logProcess.Process()
	}

	for i := 0; i < writeNum; i++ {
		go storage.Save()
	}

	m := &monitor.Monitor{
		StartTime: time.Now(),
		Data:      types.SystemInfo{},
		WebPort:   Cfg.Section("").Key("web_port").String(),
	}
	m.Start(logProcess)

}
