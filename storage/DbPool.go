/************************************************************
** @Description: storage
** @Author: george hao
** @Date:   2018-05-25 14:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-05-25 14:13
*************************************************************/
package storage

import (
	"github.com/pkg/errors"
	"sync"
	"time"
)

type PoolConfig struct {
	MinConn    int
	MaxConn    int
	Factory    func() (interface{}, error)
	Close      func()
	MaxTimeout time.Duration
}

type ChannelPool struct {
	mu      sync.Mutex
	conns   chan *Conn
	factory func() (interface{}, error)
	close   func()
	timeout time.Duration
}

type Conn struct {
	client interface{}
	t      time.Time
}

//创建连接池
func NewChannelPool(poolConfig *PoolConfig) (ChannelPool, error) {
	if poolConfig.MaxConn < 0 ||
		poolConfig.MaxConn <= 0 ||
		poolConfig.MaxConn < poolConfig.MinConn {
		return nil, errors.New("invalid settings")
	}

	c := &ChannelPool{
		conns:   make(chan *Conn, poolConfig.MaxConn),
		factory: poolConfig.Factory,
		close:   poolConfig.Close,
		timeout: poolConfig.MaxTimeout,
	}

	for i := 0; i < poolConfig.MinConn; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
		}
	}

	return c, nil

}

//获取一个链接
func (c *ChannelPool) Get() {

}

func (c *ChannelPool) Put() {

}

func (c *ChannelPool) Close() {

}

func (c *ChannelPool) Release() {

}

func (c *ChannelPool) Len() {

}
