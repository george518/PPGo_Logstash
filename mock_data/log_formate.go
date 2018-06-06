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

	Tags := `{"Path":{"id":"6","func":"url,trim_right_num"},"Method":{"id":"5","func":"trim_left_1"},"Schema":{"id":"7","func":"trim_right_1"},"Status":{"id":"8","func":""},"Ip":{"id":"0","func":""}}`
	Tag := make(map[string]interface{})
	json.Unmarshal([]byte(Tags), &Tag)

	Fields := `{"UpstreamTime":{"id":"13","func":"float64"},"RequestTime":{"id":"14","func":"float64"},"BytesSent":{"id":"9","func":"int"}}`
	Field := make(map[string]interface{})
	json.Unmarshal([]byte(Fields), &Field)

	time_key := 3
	time_format := `02/Jan/2006:15:04:05`
	time_loc := `Asia/Shanghai`

	str := `127.0.0.1 - - [31/May/2018:09:43:10 +0800] "GET /api/v0/spu/12131212 HTTP/1.0" 200 755 "-" "KeepAliveClient" "-" 0.568 0.568`
	Delimiter := ` `
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

	ret[time_key] = trim_left_1(ret[time_key])
	fmt.Println(ret[time_key])
	t, err := time.ParseInLocation(time_format, ret[time_key], loc)
	if err != nil {
		fmt.Println("error time format:", ret[time_key])
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
