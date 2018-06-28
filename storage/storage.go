/************************************************************
** @Description: storage
** @Author: george hao
** @Date:   2018-05-14 17:39
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:39
*************************************************************/
package storage

import (
	"github.com/george518/PPGo_Logstash/types"
	"github.com/influxdata/influxdb/client/v2"
	"log"
)

type Storage struct {
	Wc       chan *types.LogMessage
	Table    string
	Env      string
	WriteNum int
}

func (s *Storage) Save(pools channelPool) {

	for i := 0; i < s.WriteNum; i++ {
		go func() {
			for {
				select {
				case n := <-types.ExitSaveChan:
					log.Println(" save exit ", n)

					if len(s.Wc) > 0 {
						types.ExitSaveChan <- n
						continue
					}
					if n < s.WriteNum {
						types.ExitSaveChan <- n + 1
						continue
					}
					if n >= s.WriteNum {
						types.ExitChan <- 1
						goto EndSave
					}

				case v := <-s.Wc:
					//time.Sleep(2 * time.Second)
					go s.save(v, pools)
				}
			}

		EndSave:
			log.Println(" end save ")

		}()

	}

}

func (s *Storage) save(v *types.LogMessage, pools channelPool) {

	conn, err := pools.Get()
	if err != nil {
		log.Fatal("get conn error")
	}
	// Create a point and add to batch
	tags := v.Tags
	fields := v.Fileds
	pt, err := client.NewPoint(s.Table, tags, fields, v.TimeLocal)
	if err != nil {
		types.TypeMonitorChan <- types.TypeErrNum
		log.Println("NewPoint error:", err)
	}
	conn.bp.AddPoint(pt)

	// Write the batch

	if err := conn.cli.Write(conn.bp); err != nil {
		types.TypeMonitorChan <- types.TypeErrNum
		log.Println("InfluxDb write error:", err)
	}

	if s.Env == "development" {
		log.Println("influxdb success:", v)
		//log.Println(conn)
	}
	pools.Put(*conn)
}
