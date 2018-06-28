/************************************************************
** @Description: monitor
** @Author: george hao
** @Date:   2018-05-16 16:19
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-16 16:19
*************************************************************/
package monitor

import (
	"encoding/json"
	"github.com/george518/PPGo_Logstash/process"
	. "github.com/george518/PPGo_Logstash/types"
	"log"
	"os"
	"time"
)

type Monitor struct {
	StartTime time.Time
	Data      SystemInfo
	TpsSli    []int
	WebPort   string
}

func (m *Monitor) Start(lp *process.LogProcess) {

	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			select {
			case n := <-TypeMonitorChan:
				if n == TypeErrNum {
					m.Data.ErrNum += 1
				}
				if n == TypeHandleLine {
					m.Data.HandleLine += 1
				}

			case <-ticker.C:
				m.TpsSli = append(m.TpsSli, m.Data.HandleLine)
				if len(m.TpsSli) > 2 {
					m.TpsSli = m.TpsSli[1:]
				}
				m.PrintMonitor(lp)
			}
		}

	}()

	//http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
	//	m.Data.RunTime = time.Now().Sub(m.StartTime).String()
	//	m.Data.ReadChanLen = len(lp.Rc)
	//	m.Data.WriteChanLen = len(lp.Wc)
	//
	//	if len(m.TpsSli) >= 2 {
	//		m.Data.Tps = float64(m.TpsSli[1]-m.TpsSli[0]) / 5
	//	}
	//
	//	ret, err := json.MarshalIndent(m.Data, "", "\t")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	io.WriteString(writer, string(ret))
	//})
	//
	//http.ListenAndServe(":"+m.WebPort, nil)
}

func (m *Monitor) PrintMonitor(lp *process.LogProcess) {
	m.Data.RunTime = time.Now().Sub(m.StartTime).String()
	m.Data.ReadChanLen = len(lp.Rc)
	m.Data.WriteChanLen = len(lp.Wc)

	if len(m.TpsSli) >= 2 {
		m.Data.Tps = float64(m.TpsSli[1]-m.TpsSli[0]) / 5
	}

	ret, err := json.MarshalIndent(m.Data, "", "\t")
	if err != nil {
		log.Println(err)
	}

	//监控日志记录

	file, err := os.OpenFile(MonitorLogPath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println(" monitor log open error", err)
	}
	_, err = file.Write(ret)
	if err != nil {
		log.Println(" monitor log write error", err)
	}
}
