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
	"io"
	"log"
	"net/http"
	"time"
)

type Monitor struct {
	StartTime time.Time
	Data      SystemInfo
	TpsSli    []int
	WebPort   string
}

func (m *Monitor) Start(lp *process.LogProcess) {
	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.Data.ErrNum += 1
			case TypeHandleLine:
				m.Data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-ticker.C
			m.TpsSli = append(m.TpsSli, m.Data.HandleLine)
			if len(m.TpsSli) > 2 {
				m.TpsSli = m.TpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
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
		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":"+m.WebPort, nil)
}
