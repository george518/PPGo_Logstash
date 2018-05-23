/************************************************************
** @Description: storage
** @Author: george hao
** @Date:   2018-05-14 17:39
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:39
*************************************************************/
package storage

import (
	"github.com/george518/PPGo_Logstash/process"
	"github.com/george518/PPGo_Logstash/types"
	"github.com/influxdata/influxdb/client/v2"
	"log"
)

type Storage struct {
	Wc                                                chan *process.Message
	DbUrl, DbPort, DbUser, DbPwd, DbName, DbPrecision string
	DbTable                                           string
	Env                                               string
}

func (s *Storage) Save() {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     s.DbUrl + ":" + s.DbPort,
		Username: s.DbUser,
		Password: s.DbPwd,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.DbName,
		Precision: s.DbPrecision,
	})
	if err != nil {
		log.Fatal(err)
	}

	for v := range s.Wc {
		// Create a point and add to batch
		tags := map[string]string{
			"Path":   v.Path,
			"Method": v.Method,
			"Schema": v.Scheme,
			"Status": v.Status,
			"Ip":     v.Ip,
		}
		//tags := map[string]string{"cpu": "cpu-total"}
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		pt, err := client.NewPoint(s.DbTable, tags, fields, v.TimeLocal)
		if err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println(err)
			continue
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := c.Write(bp); err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println(err)
			continue
		}

		if s.Env == "development" {
			log.Println("influxdb success:", v)
		}

	}
}
