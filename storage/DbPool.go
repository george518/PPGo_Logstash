/************************************************************
** @Description: storage
** @Author: george hao
** @Date:   2018-05-25 14:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-25 14:13
*************************************************************/
package storage

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/george518/PPGo_Logstash/config"
	"github.com/influxdata/influxdb/client/v2"
)

var ErrClosed = errors.New("pool is closed")

//PoolConfig 连接池相关配置
type PoolConfig struct {
	//连接池中拥有的最小连接数
	InitialCap int
	//连接池中拥有的最大的连接数
	MaxCap int
	//生成连接的方法
	Factory func() (InfluxDb, error)
	//关闭链接的方法
	Close func(db InfluxDb) error
	//链接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

//channelPool 存放链接信息
type channelPool struct {
	mu          sync.Mutex
	conns       chan *idleConn
	factory     func() (InfluxDb, error)
	close       func(db InfluxDb) error
	idleTimeout time.Duration
}
type InfluxDb struct {
	cli client.Client
	bp  client.BatchPoints
}

type idleConn struct {
	conn InfluxDb
	t    time.Time
}

//NewChannelPool 初始化链接
func NewChannelPool(poolConfig *PoolConfig) (*channelPool, error) {
	if poolConfig.InitialCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.InitialCap > poolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &channelPool{
		conns:       make(chan *idleConn, poolConfig.MaxCap),
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
	}

	for i := 0; i < poolConfig.MaxCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	return c, nil
}

//getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

//Get 从pool中取一个连接
func (c *channelPool) Get() (ifd *InfluxDb, err error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					c.Close(wrapConn.conn)
					continue
				}
			}
			return &wrapConn.conn, nil
		default:
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}
			return &conn, nil
		}
	}
}

//Put 将连接放回pool中
func (c *channelPool) Put(conn InfluxDb) error {
	if &conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		return c.Close(conn)
	}

	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return c.Close(conn)
	}
}

//Close 关闭单条连接
func (c *channelPool) Close(conn InfluxDb) error {
	if &conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.close(conn)
}

//Release 释放连接池中所有链接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

//Len 连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}

func DbFactory() (ifd InfluxDb, err error) {
	s := config.LoadConfig()
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     s.Storage.Key("Url").String() + ":" + s.Storage.Key("Port").String(),
		Username: s.Storage.Key("User").String(),
		Password: s.Storage.Key("Pwd").String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.Storage.Key("Name").String(),
		Precision: s.Storage.Key("Precision").String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	ifd.cli = c
	ifd.bp = bp
	return ifd, err
}

func DbClose(db InfluxDb) error {
	return db.cli.Close()
}
