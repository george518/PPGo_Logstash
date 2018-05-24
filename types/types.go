/************************************************************
** @Description: types
** @Author: george hao
** @Date:   2018-05-16 17:28
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-16 17:28
*************************************************************/
package types

import "time"

// 系统状态监控
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 总处理日志行数
	Tps          float64 `json:"tps"`          // 系统吞出量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel 长度
	RunTime      string  `json:"runTime"`      // 运行总时间
	ErrNum       int     `json:"errNum"`       // 错误数
}

//信息格式
type Message struct {
	TimeLocal                        time.Time
	BytesSent                        int
	Path, Method, Scheme, Status, Ip string
	UpstreamTime, RequestTime        float64
}

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

var TypeMonitorChan = make(chan int, 200)

type Conf struct {
	AppMode    string
	WebPort    string
	ReadNum    int
	ProcessNum int
	WriteNum   int
	StorageDb  StorageDb
	LogInfo    LogInfo
}

type StorageDb struct {
	Type      string
	Url       string
	Port      string
	User      string
	Pwd       string
	Name      string
	Precision string
	Table     string
}

type LogInfo struct {
	Path    string
	Regexp  string
	Time    string
	TimeLoc string
}
