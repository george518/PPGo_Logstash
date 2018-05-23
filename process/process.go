/************************************************************
** @Description: process
** @Author: george hao
** @Date:   2018-05-14 17:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:13
*************************************************************/
package process

import (
	"github.com/george518/PPGo_Logstash/types"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	TimeLocal                        time.Time
	BytesSent                        int
	Path, Method, Scheme, Status, Ip string
	UpstreamTime, RequestTime        float64
}
type LogProcess struct {
	Wc         chan *Message
	Rc         chan []byte
	Regexp     string
	TimeLoc    string
	TimeFormat string
}

func (lp *LogProcess) Process() {

	r := regexp.MustCompile(lp.Regexp)

	for v := range lp.Rc {

		//fmt.Println(string(v))
		ret := r.FindStringSubmatch(string(v))
		loc, _ := time.LoadLocation(lp.TimeLoc)
		if len(ret) != 13 {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation(lp.TimeFormat, ret[4], loc)
		if err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("ParseInLocation ", ret[4], err)
			continue
		}

		message.TimeLocal = t

		bytesSent, _ := strconv.Atoi(ret[7])
		message.BytesSent = bytesSent

		reqSli := strings.Split(ret[5], " ")
		if len(reqSli) != 3 {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("strings.Split", ret[5])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			types.TypeMonitorChan <- types.TypeErrNum
			log.Println("url parse fail:", reqSli[1])
			continue
		}

		// /api/v0/spu/12123212=>/api/v0/spu

		message.Path = u.Path
		if n := strings.Count(u.Path, "/"); n == 4 {
			pathByte := []byte(u.Path)
			pathByte = pathByte[0:strings.LastIndex(u.Path, "/")]
			message.Path = string(pathByte)
		}

		message.Scheme = reqSli[2]
		message.Status = ret[6]

		//message.Ip = ret[1]

		upstreamTime, _ := strconv.ParseFloat(ret[11], 64)
		requestTime, _ := strconv.ParseFloat(ret[12], 64)

		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		lp.Wc <- message
	}
}
