package registry

import (
	"github.com/go-zk/zk"
	"sync"
	"time"
)

//初始化
func (z *ZkRegistry) init() {
	z.connect()
	//初始化根路径
	z.ensureRoot()
	//watch 监听服务
	z.watchListens()
}

//zk链接
func (z *ZkRegistry) connect() {
	conn, _, err := zk.Connect(z.opts.hosts, time.Duration(z.opts.timeout)*time.Second)
	if err != nil {
		panic(err.Error())
	}
	z.conn = conn
}

//初始化根路径
func (z *ZkRegistry) ensureRoot() error {
	exists, _, err := z.conn.Exists(z.opts.prefix)
	if err != nil {
		return err
	}
	if !exists {
		_, err := z.conn.Create(z.opts.prefix, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//watch 所有监听服务首次监听需阻塞
func (z *ZkRegistry) watchListens() {
	var wg = sync.WaitGroup{}
	for _, serviceName := range z.opts.listens {
		wg.Add(1)
		go z.watchNode(serviceName, &wg)
	}
	wg.Wait()
}
