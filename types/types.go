/************************************************************
** @Description: types
** @Author: george hao
** @Date:   2018-05-16 17:28
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-16 17:28
*************************************************************/
package types

import (
	"time"
)

// 系统状态监控
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 总处理日志行数
	Tps          float64 `json:"tps"`          // 系统吞出量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel 长度0
	RunTime      string  `json:"runTime"`      // 运行总时间
	ErrNum       int     `json:"errNum"`       // 错误数
}

type LogMessage struct {
	TimeLocal time.Time
	Tags      map[string]string
	Fileds    map[string]interface{}
}

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

var TypeMonitorChan = make(chan int, 200)

var ExitReadChan = make(chan int)
var ExitProcessChan = make(chan int)
var ExitSaveChan = make(chan int)
var ExitChan = make(chan int)

var MonitorLogPath string
