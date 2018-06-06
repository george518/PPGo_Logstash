/************************************************************
** @Description: process
** @Author: george hao
** @Date:   2018-05-14 17:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:13
*************************************************************/
package process

import (
	"encoding/json"
	"github.com/george518/PPGo_Logstash/types"
	"github.com/go-ini/ini"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogProcess struct {
	Wc      chan *types.LogMessage
	Rc      chan []byte
	LogInfo *ini.Section
}

func (lp *LogProcess) Process() {
	logType := lp.LogInfo.Key("Type").String()
	switch logType {
	case "nginx_log":
		nginx_process(lp)
	default:
		log.Fatalln("illegal logtype :", logType)
	}

}

func nginx_process(lp *LogProcess) {
	loc, err := time.LoadLocation(lp.LogInfo.Key("TimeLoc").String())
	if err != nil {
		log.Fatal("Local time formatting error")
	}
	for v := range lp.Rc {
		delimiter := lp.LogInfo.Key("Delimiter").String()
		if delimiter == "" {
			delimiter = " "
		}
		ret := strings.Split(string(v), delimiter)

		//处理Tags
		Tags := lp.LogInfo.Key("Tags").String()
		Tag := make(map[string]interface{})
		err := json.Unmarshal([]byte(Tags), &Tag)

		if err != nil {
			log.Println("tags json unmarshal error")
			continue
		}

		tags := make(map[string]string)
		for k, v := range Tag {
			iv := v.(map[string]interface{})
			id := -1
			// 查找键值是否存在iv
			if isId, ok := iv["id"]; ok {
				idInt, err := strconv.Atoi(isId.(string))
				if err != nil {
					log.Println("strconv.Atoi id conv err", iv["id"])
					continue
				}
				id = idInt
			} else {
				log.Println("index id is not exist")
				continue
			}

			if id == -1 {
				log.Println("id illegal -1")
				continue
			}

			f := iv["func"].(string)

			if f == "" {
				tags[k] = ret[id]
				continue
			}

			tags[k] = ret[id]
			fs := strings.Split(f, ",")

			for _, vf := range fs {
				switch vf {
				case "url":
					tags[k] = url_format(tags[k])
				case "trim_right_num":
					tags[k] = trim_right_num(tags[k])
				case "trim_left_1":
					tags[k] = trim_left_1(tags[k])
				case "trim_right_1":
					tags[k] = trim_right_1(tags[k])
				default:
					tags[k] = ret[id]
				}
			}

		}

		//处理Fileds
		Fields := lp.LogInfo.Key("Fields").String()
		Field := make(map[string]interface{})
		err = json.Unmarshal([]byte(Fields), &Field)

		if err != nil {
			log.Println("fields json unmarshal error")
			continue
		}

		fields := make(map[string]interface{})
		for k, v := range Field {
			iv := v.(map[string]interface{})
			id := -1
			// 查找键值是否存在iv
			if isId, ok := iv["id"]; ok {

				idInt, err := strconv.Atoi(isId.(string))
				if err != nil {
					log.Println("fields strconv.Atoi id conv err", iv["id"])
					continue
				}
				id = idInt
			} else {
				log.Println("fields id is not exist")
				continue
			}

			if id == -1 {
				log.Println("fields id illegal -1")
				continue
			}

			f := iv["func"].(string)

			if f == "" {
				fields[k] = ret[id] //strconv.ParseFloat(ret[id], 64)
				continue
			}
			fs := strings.Split(f, ",")

			for _, vf := range fs {
				switch vf {
				case "float64":
					fields[k], _ = strconv.ParseFloat(ret[id], 64)

				case "float32":
					fields[k], _ = strconv.ParseFloat(ret[id], 32)
				case "int":
					fields[k], _ = strconv.Atoi(ret[id])

				default:
					fields[k] = 0
				}
			}
		}

		//处理时间 `{"Time":{"id":"3","func":"trim_left_1"}}`
		Times := lp.LogInfo.Key("Times").String()
		Time := make(map[string]string)
		json.Unmarshal([]byte(Times), &Time)
		id := -1
		if idstr, ok := Time["id"]; ok {
			id, _ = strconv.Atoi(idstr)
		} else {
			log.Println("time strconv.Atoi id is not exist")
			continue
		}

		if id == -1 {
			log.Println("time id illegal:", id)
			continue
		}
		timeStr := ret[id]
		if funcstr, ok := Time["func"]; ok {
			if funcstr != "" {
				switch funcstr {
				case "trim_left_1":
					timeStr = trim_left_1(timeStr)
				case "trim_right_1":
					timeStr = trim_right_1(timeStr)
				}
			}
		}

		t, err := time.ParseInLocation(lp.LogInfo.Key("TimeFormat").String(), timeStr, loc)
		if err != nil {
			log.Println(" time is error")
			continue
		}

		logmessage := &types.LogMessage{
			TimeLocal: t,
			Tags:      tags,
			Fileds:    fields,
		}

		lp.Wc <- logmessage
	}
}

func url_format(url_str string) string {
	u, err := url.Parse(url_str)
	if err != nil {
		return url_str
	}
	return u.Path
}

func trim_right_num(str string) string {
	pat := `/[0-9].*`
	re, err := regexp.Compile(pat)
	if err != nil {
		return str
	}
	return re.ReplaceAllString(str, "")
}

func trim_left_1(str string) string {
	bt := []byte(str)
	return string(bt[1:])
}

func trim_right_1(str string) string {
	bt := []byte(str)
	return string(bt[:len(bt)-1])
}
