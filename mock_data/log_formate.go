/************************************************************
** @Description: 测试日志格式
** @Author: george hao
** @Date:   2018-06-05 09:56
** @Last Modified by:  george hao
** @Last Modified time: 2018-06-05 09:56
*************************************************************/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	str := `GET /api/v0/order?sign=c96be7781f8579df8f52edb670e0c316&page_size=50&ts=1528353264&end_time=2018-06-07 14:34:24&start_time=2018-06-07 14:05:24&method=GET&app_key=0137&page_no=1 HTTP/1.1`
	s := split_str(str, " ", -1)
	log.Println(s)

	arr := []int{1, 2, 3}

	if t, ok := arr[3]; ok {
		log.Println("no")
	} else {
		log.Println("yes")
	}

}

func parse() {

	Tags := `{"Path":{"id":"4","func":"split_str,url,trim_right_num","c_id":"1"},"Method":{"id":"4","func":"split_str,trim_left_1","c_id":"0"},"Schema":{"id":"4","func":"split_str,trim_right_1","c_id":"-1"},"Status":{"id":"5","func":""},"Ip":{"id":"0","func":""}}`
	Tag := make(map[string]interface{})
	json.Unmarshal([]byte(Tags), &Tag)

	Fields := `{"UpstreamTime":{"id":"10","func":"float64"},"RequestTime":{"id":"11","func":"float64"},"BytesSent":{"id":"6","func":"int"}}`
	Field := make(map[string]interface{})
	json.Unmarshal([]byte(Fields), &Field)

	time_format := `02/Jan/2006:15:04:05 +0800`
	time_loc := `Asia/Shanghai`

	Timesf := `{"id":"3","func":"trim_left_1,trim_right_1"}`

	//str := `127.0.0.1||-||-||[31/May/2018:09:43:10 +0800]||"GET /api/v0/spu/12131212&time=2017-12-18 12:00:00&end_time=2017-12-30 12:00:00&ts=12  HTTP/1.0"||200||755||"-"||"KeepAliveClient"||"-"||0.568||0.568`
	str := `140.207.52.210||-||-||[07/Jun/2018:14:34:25 +0800]||"GET /api/v0/order?sign=c96be7781f8579df8f52edb670e0c316&page_size=50&ts=1528353264&end_time=2018-06-07 14:34:24&start_time=2018-06-07 14:05:24&method=GET&app_key=0137&page_no=1 HTTP/1.1"||200||72||"-"||"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1;SV1)"||"-"||0.169||0.169`
	Delimiter := `||`
	ret := strings.Split(str, Delimiter)

	for ks, vs := range ret {
		fmt.Println(ks, "=>", vs)
	}

	tags := make(map[string]string)
	for k, v := range Tag {
		iv := v.(map[string]interface{})
		id := -1
		// 查找键值是否存在iv
		if isId, ok := iv["id"]; ok {

			idInt, err := strconv.Atoi(isId.(string))
			if err != nil {
				log.Println("id conv err", iv["id"])
				continue
			}
			id = idInt
		} else {
			log.Println("id is not exist")
			continue
		}

		if id == -1 {
			log.Println("id illegal")
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
			case "split_str":
				cid, _ := strconv.Atoi(iv["c_id"].(string))
				tags[k] = split_str(tags[k], " ", cid)
			case "url":
				tags[k] = url_format(tags[k])
			case "trim_right_num":
				tags[k] = trim_right_num(tags[k])
			case "trim_left_1":
				tags[k] = trim_left_1(tags[k])
			case "trim_right_1":
				tags[k] = trim_right_1(tags[k])
			default:
				tags[k] = vf
			}
		}

	}

	fields := make(map[string]interface{})
	for k, v := range Field {
		iv := v.(map[string]interface{})
		id := -1
		// 查找键值是否存在iv
		if isId, ok := iv["id"]; ok {

			idInt, err := strconv.Atoi(isId.(string))
			if err != nil {
				log.Println("id conv err", iv["id"])
				continue
			}
			id = idInt
		} else {
			log.Println("id is not exist")
			continue
		}

		if id == -1 {
			log.Println("id illegal")
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
			case "int":
				fields[k], _ = strconv.Atoi(ret[id])
			default:
				fields[k] = 0
			}
		}
	}

	loc, _ := time.LoadLocation(time_loc)

	//处理时间
	Times := Timesf
	Time := make(map[string]string)
	json.Unmarshal([]byte(Times), &Time)

	id := -1
	if idstr, ok := Time["id"]; ok {
		id, _ = strconv.Atoi(idstr)
	} else {
		log.Println("time strconv.Atoi id is not exist")
	}

	if id == -1 {
		log.Println("time id illegal:", id)

	}
	timeStr := ret[id]

	if funcstr, ok := Time["func"]; ok {
		fs := strings.Split(funcstr, ",")
		for _, vf := range fs {
			switch vf {
			case "trim_left_1":
				timeStr = trim_left_1(timeStr)
			case "trim_right_1":
				timeStr = trim_right_1(timeStr)
			}
		}

	}

	t, err := time.ParseInLocation(time_format, timeStr, loc)
	if err != nil {
		fmt.Println("error time format:", timeStr)
	}

	//fields := map[string]interface{}{
	//	"UpstreamTime": v.UpstreamTime,
	//	"RequestTime":  v.RequestTime,
	//	"BytesSent":    v.BytesSent,
	//}

	fmt.Println(tags)
	fmt.Println(fields)

	fmt.Println(t)

}

func replace_str(str string) string {
	return strings.Replace(str, "/", "|", 1)
}

func split_str(str, delimiter string, id int) string {
	arr := strings.Split(str, delimiter)
	if id == -1 {
		id = len(arr) - 1
	}
	if id < len(arr) {
		return arr[id]
	}
	return str
}

func url_format(url_str string) string {
	u, err := url.Parse(url_str)
	if err != nil {
		return "err"
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

func print(ret []string) {
	for k, v := range ret {
		fmt.Println(k, "=>", v)
	}
}

func ba() {
	str := `127.0.0.1 - - [31/May/2018:09:43:10 +0800] "GET /api/v0/spu HTTP/1.0" 200 755 "-" "KeepAliveClient" "-" 0.568 0.568`
	Regexp := `([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`

	//Time    = 02/Jan/2006:15:04:05 +0800=
	t := time.Now()
	r := regexp.MustCompile(Regexp)
	ret := r.FindStringSubmatch(str)
	fmt.Println(time.Since(t))
	print(ret)

	t1 := time.Now()

	Delimiter := " "

	//Delimiter := " "
	ret1 := strings.Split(str, Delimiter)

	fmt.Println(time.Since(t1))
	print(ret1)
}
