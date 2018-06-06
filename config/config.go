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

//func LoadConfig() *types.Conf {
//	cfg, err := ini.Load("./config/conf.ini")
//	conf := new(types.Conf)
//	err = cfg.MapTo(conf)
//
//	if err != nil {
//		log.Fatalf("Fail to read file: %v", err)
//	}
//	return conf
//}

type Conf struct {
	Global  *ini.Section
	Storage *ini.Section
	LogType *ini.Section
}

func LoadConfig(config_file string) *Conf {
	cfg, err := ini.Load(config_file)

	if err != nil {
		log.Fatalf("Fail to read config file:%v", err)
	}

	globalConf := cfg.Section("")
	storageType := globalConf.Key("StorageType").String()
	logType := globalConf.Key("LogType").String()
	storageConf := cfg.Section(storageType)
	logConf := cfg.Section(logType)

	return &Conf{
		Global:  globalConf,
		Storage: storageConf,
		LogType: logConf,
	}
}
