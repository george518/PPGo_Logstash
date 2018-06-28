/************************************************************
** @Description: logdig
** @Author: george hao
** @Date:   2018-05-14 17:04
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-14 17:04
*************************************************************/
package logdig

import (
	"bufio"
	"github.com/george518/PPGo_Logstash/types"
	"github.com/go-ini/ini"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type LogData struct {
	Rc   chan []byte
	Conf *ini.Section
}

func (ld *LogData) Read() {

	paths := strings.Split(ld.Conf.Key("Path").String(), ",")

	fileNum := len(paths)
	go func() {
		for _, path := range paths {
			go ReadFile(path, ld.Rc, fileNum)
		}
	}()

}

func ReadFile(path string, ch chan []byte, fileNum int) {
	//TODO 如果文件资源变化，需要重新打开文件句柄
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		select {
		case n := <-types.ExitReadChan:
			types.ExitReadChan <- n + 1
			f.Close()
			log.Println(" Stop reading the file,ID:", path, n)
			if n == fileNum {
				types.ExitProcessChan <- 1
				goto ForEnd
			}

		default:
			line, err := rd.ReadBytes('\n')
			if err == io.EOF {
				time.Sleep(500 * time.Microsecond)
				continue
			} else if err != nil {
				log.Println(err)
				continue
			}
			types.TypeMonitorChan <- types.TypeHandleLine
			ch <- line[:len(line)-1]

		}

	}

ForEnd:
	log.Println(" Completely stop reading", path)
}
