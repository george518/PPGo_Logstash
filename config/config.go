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

type Config struct {
}

func Load() ini.File {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	return *cfg
}
