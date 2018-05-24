/************************************************************
** @Description: config
** @Author: george hao
** @Date:   2018-05-16 08:46
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-16 08:46
*************************************************************/
package config

import (
	"github.com/george518/PPGo_Logstash/types"
	"github.com/go-ini/ini"
	"log"
)

func LoadConfig() *types.Conf {
	cfg, err := ini.Load("./config/conf.ini")
	conf := new(types.Conf)
	err = cfg.MapTo(conf)

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	return conf
}
