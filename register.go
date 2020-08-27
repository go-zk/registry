package registry

import (
	"errors"
	"math/rand"
)

//随机获取一个服务节点
func (z *ZkRegistry) GetRandNode(serviceName string) (*Node, error) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	total := len(z.register[serviceName])
	//fmt.Println(total)
	if total < 1 {
		return nil, errors.New("service " + serviceName + " not node")
	}
	var data = make([]*Node, 0)
	for _, node := range z.register[serviceName] {
		data = append(data, node)
	}
	return data[rand.Intn(total)], nil
}

//获取服务所有节点(外部负载均衡使用)
func (z *ZkRegistry) GetAllNode(serviceName string) []*Node {
	z.mu.RLock()
	defer z.mu.RUnlock()
	var data = make([]*Node, 0)
	for _, node := range z.register[serviceName] {
		data = append(data, node)
	}
	return data
}

//获取所有服务名称
func (z *ZkRegistry) GetAllService() []string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	var data = make([]string, 0)
	for serviceName, _ := range z.register {
		data = append(data, serviceName)
	}
	return data
}

//添加节点到本地map
func (z *ZkRegistry) setServiceNode(serviceName string, node *Node) error {
	z.mu.Lock()
	defer z.mu.Unlock()
	if node.Id == "" {
		return errors.New("添加节点失败，没有节点ID")
	}
	if _, ok := z.register[serviceName]; !ok {
		z.register[serviceName] = make(map[string]*Node)
	}
	z.register[serviceName][node.Id] = node
	return nil
}

//删除节点
func (z *ZkRegistry) DelServiceNode(serviceName string, node *Node) error {
	z.mu.Lock()
	defer z.mu.Unlock()
	if _, ok := z.register[serviceName][node.Id]; !ok {
		return errors.New("node not found :" + node.Id)
	}
	delete(z.register[serviceName], node.Id)
	return nil
}
