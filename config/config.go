/************************************************************
** @Description: config
** @Author: george hao
** @Date:   2018-05-16 08:46
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-16 08:46
*************************************************************/
package config

import (
	"github.com/go-ini/ini"
	"log"
)

//app_mode = development
//
//time_loc = Asia/Shanghai
//web_port = 8080
//
//read_num    = 1
//process_num = 5
//write_num   = 50
//
//[storage]
//db_type      = influxDb
//db_url       = http://10.32.33.27
//db_port      = 8086
//db_user      = george
//db_pwd       = 123456
//db_name      = georgetest
//db_precision = s
//db_table     = nginx_log
//
//[log_info]
//log_path    = ./testData/access.log
//log_regexp  = `([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`
//log_time    = 02/Jan/2006:15:04:05 +0800
//
//#log format
//#      log_format  access  '$remote_addr - $remote_user [$time_local] "$request" '
//#                          '$status $body_bytes_sent "$http_referer" '
//#                         '"$http_user_agent" "$http_x_forwarded_for" '
//#                          '$upstream_response_time $request_time';

type ConfigGlobal struct {
	AppMode    string
	TimeLoc    string
	WebPort    string
	ReadNum    int
	ProcessNum int
	WriteNum   int
	StorageDb
	LogInfo
}

type StorageDb struct {
	Type      string
	Url       string
	Port      string
	User      string
	Pwd       string
	Name      string
	Precision string
	Table     string
}

type LogInfo struct {
	Path       string
	Regexp     string
	TimeFormat string
}

func Load() ini.File {
	cfg, err := ini.Load("./config/config.ini")
	glob := new(ConfigGlobal)
	//err = cfg.MapTo(glob)
	cfg2, err := ini.Load("./config/conf.ini")
	err = cfg2.MapTo(glob)

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	return *cfg
}

func LoadConfig() *ConfigGlobal {
	cfg, err := ini.Load("./config/conf.ini")
	glob := new(ConfigGlobal)
	err = cfg.MapTo(glob)

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	return glob
}
