package registry

import (
	"encoding/json"
	"fmt"
	"github.com/go-basic/uuid"
	"github.com/go-zk/zk"
	"log"
	"math/rand"
	"sync"
	"time"
)

//注册中心
type ZkRegistry struct {
	opts Options
	conn *zk.Conn
	//本地注册表
	mu       sync.RWMutex
	register map[string]map[string]*Node
}

//创建注册中心
func NewZkRegistry(opts ...Option) (*ZkRegistry, error) {
	z := &ZkRegistry{
		register: make(map[string]map[string]*Node),
	}
	options := newOptions(opts...)
	// set opts
	z.opts = options
	//初始化随机
	rand.Seed(time.Now().UnixNano())
	//服务初始化
	z.init()
	return z, nil
}

//服务注册(注册节点到注册中心)
func (z *ZkRegistry) Register(serviceName string, node *Node) (err error) {
	path := z.getPath(serviceName)
	ex, _, err := z.conn.Exists(path)
	if err != nil {
		log.Println("service path:", path)
		return err
	}
	if !ex {
		//持久化服务节点
		_, err = z.conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Println("Create service path error", path)
			return err
		}
	}
	//node节点（临时节点）
	nodePath := path + "/" + node.Id
	ex, _, err = z.conn.Exists(nodePath)
	if err != nil {
		log.Println("node path error", nodePath)
		return err
	}
	if !ex {
		addJava(node)
		node.Name = serviceName //将服务名称加入节点
		data, _ := json.Marshal(node)
		_, err = z.conn.Create(nodePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Println("Create node path error", nodePath)
			return err
		}
	}
	return
}

//服务反注册
func (z *ZkRegistry) Deregister(serviceName string, node *Node) (err error) {
	path := z.getPath(serviceName)
	nodePath := path + "/" + node.Id
	ex, stat, err := z.conn.Exists(nodePath)
	if err != nil {
		log.Println("node path error", nodePath)
		return err
	}
	if ex {
		if err = z.conn.Delete(nodePath, stat.Version); err != nil {
			return err
		}
	}
	return nil
}

//手动添加服务节点(不会注册到注册中心)
func (z *ZkRegistry) SetServiceNode(serviceName, address string) error {
	node := &Node{
		Id:      uuid.New(),
		Name:    serviceName,
		Address: address,
	}
	return z.setServiceNode(serviceName, node)
}

//随机获取一个服务节点
func (z *ZkRegistry) GetServerNode(serviceName string) (*Node, error) {
	return z.GetRandNode(serviceName)
}

//关闭服务
func (z *ZkRegistry) Close() {
	z.conn.Close()
	log.Println("zk conn close")
	return
}

//获取服务列表
func (z *ZkRegistry) getServices() (services []*Service, err error) {
	services = make([]*Service, 0)
	path := z.getPath("")
	list, _, err := z.conn.Children(path)
	if err != nil {
		return
	}
	for _, item := range list {
		var service = &Service{}
		service.Name = item
		services = append(services, service)
	}
	return
}

//获取服务节点列表
func (z *ZkRegistry) getServiceNodes(serviceName string) (nodes []*Node, err error) {
	nodes = make([]*Node, 0)
	path := z.getPath(serviceName)
	list, _, err := z.conn.Children(path)
	if err != nil {
		return
	}
	for _, item := range list {
		itemData, _, _ := z.conn.Get(path + "/" + item)
		var node Node
		json.Unmarshal(itemData, &node)
		//fmt.Println(fmt.Sprintf("%+v", node))
		nodes = append(nodes, &node)
	}
	return
}

//watch机制，监听节点变化(service节点&node节点)
func (z *ZkRegistry) watchNode(serviceName string, wg *sync.WaitGroup) {
	path := z.getPath(serviceName)
	conn := z.conn
	go func(wg *sync.WaitGroup) {
		var i int
		for {
			nodeIds, _, events, err := conn.ChildrenW(path)
			if err != nil {
				log.Println("watch[" + serviceName + "]err:" + err.Error())
				time.Sleep(time.Duration(z.opts.timeout) * time.Second)
				continue
			}
			//获取信息
			var newNodes = make(map[string]bool)
			for _, nodeId := range nodeIds {
				newData, _, err := conn.Get(path + "/" + nodeId)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				var node Node
				err = json.Unmarshal(newData, &node)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				newNodes[nodeId] = true
				if node.Id == "" { //兼容node没有定义ID的情况
					node.Id = nodeId
				}
				err = z.setServiceNode(serviceName, &node)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				log.Println(fmt.Sprintf("%+v", &node))
			}
			//不在注册中心的节点删除
			for id, node := range z.register[serviceName] {
				if _, ok := newNodes[id]; !ok {
					err = z.DelServiceNode(serviceName, node)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}
			if i == 0 {
				wg.Done()
			}
			i = 1
			select {
			case evt := <-events:
				if evt.Err != nil {
					log.Println(evt.Err.Error())
				}
				log.Printf("ChildrenW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}(wg)
}

//获取服务路径
func (z *ZkRegistry) getPath(serviceName string) string {
	return z.opts.prefix + "/" + serviceName
}
