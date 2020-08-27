package main

import (
	"github.com/go-zk/registry"
	"log"
)

//服务发现示例
func main() {
	zkRegistry, err := registry.NewZkRegistry(
		registry.Hosts([]string{"127.0.0.1:2181"}),
		registry.Prefix("/zk-registry"),
		registry.Timeout(15),
		registry.Listens([]string{"test-service"}),
	)
	if err != nil {
		panic(err.Error())
	}
	defer zkRegistry.Close()
	//随机负载均衡获取一个节点
	node, err := zkRegistry.GetServerNode("test-service")
	if err != nil {
		log.Println(err)
	}
	log.Println("随机负载均衡获取的node:", node)

	//外部负载均衡时使用
	nodes := zkRegistry.GetAllNode("test-service")
	log.Println("test-service下的节点列表")
	for _, node := range nodes {
		log.Println(node)
	}
}
