/************************************************************
** @Description: testData
** @Author: george hao
** @Date:   2018-05-15 11:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-15 11:13
*************************************************************/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	/**
	127.0.0.1 - - [17/May/2018:13:30:01 +0800] "GET /api/v0/product HTTP/1.1" 200 119 "-" "-" - "0.065"
	*/
	const path = "./access.log"

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Open file err: %s", err.Error()))
	}
	defer file.Close()

	for {

		for i := 1; i < 4; i++ {
			now := time.Now()
			rand.Seed(now.UnixNano())
			paths := []string{"/api/v0/product", "/api/v0/order", "/api/v0/spu", "/api/v0/sku", "/api/v0/image", "/api/v0/m_sku", "/api/v0/m_spu", "/api/v0/m_image"}
			path := paths[rand.Intn(len(paths))]
			requestTime := rand.Float64()
			if path == "/api/v0/product" {
				requestTime = requestTime + 1.4
			}

			method := "POST"
			if now.UnixNano()/1000%2 == 1 {
				method = "GET"
			}
			dateTime := now.Format("02/Jan/2006:15:04:05")
			code := 200
			if now.Unix()%10 == 1 {
				code = 500
			}
			bytesSend := rand.Intn(1000) + 500
			if path == "/api/v0/product" {
				bytesSend = bytesSend + 1000
			}

			upstreamResponseTime := requestTime
			if now.Unix()%10 == 1 {
				upstreamResponseTime = requestTime + 0.001
			}
			line := fmt.Sprintf("127.0.0.1 - - [%s +0800] \"%s %s HTTP/1.0\" %d %d \"-\" \"KeepAliveClient\" \"-\" %.3f %.3f\n", dateTime, method, path, code, bytesSend, upstreamResponseTime, requestTime)

			_, err := file.Write([]byte(line))

			fmt.Println(line)
			if err != nil {
				log.Println("writeToFile error:", err)
			}
		}

		time.Sleep(4 * time.Second)
		//time.Sleep(time.Millisecond * 200)

		time.Sleep(2 * time.Second)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
