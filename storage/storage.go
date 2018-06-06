/************************************************************
** @Description: storage
** @Author: george hao
** @Date:   2018-05-14 17:39
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:39
*************************************************************/
package storage

import (
	"log"

	"github.com/george518/PPGo_Logstash/types"
	"github.com/influxdata/influxdb/client/v2"
)

type Storage struct {
	Wc    chan *types.LogMessage
	Table string
	Env   string
}

func (s *Storage) Save(pools channelPool) {

	for v := range s.Wc {
		conn, err := pools.Get()
		if err != nil {
			log.Fatal("get conn error")
		}
		// Create a point and add to batch
		tags := v.Tags
		//tags := map[string]string{"cpu": "cpu-total"}
		fields := v.Fileds
		pt, err := client.NewPoint(s.Table, tags, fields, v.TimeLocal)
		if err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("NewPoint error:", err)
			continue
		}
		conn.bp.AddPoint(pt)

		// Write the batch

		if err := conn.cli.Write(conn.bp); err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("InfluxDb write error:", err)
			continue
		}

		if s.Env == "development" {
			log.Println("influxdb success:", v)
			//log.Println(conn)
		}
		pools.Put(*conn)
	}
}
