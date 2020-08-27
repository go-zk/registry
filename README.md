# zk注册中心
实现了服务注册与发现逻辑，模块独立，引入方便。

支持：

- 服务注册
- 服务发现
- 可添加node到本地注册表，不注册到zk
- 支持随机负载均衡
- 支持外部负载均衡
- spring cloud（java）服务发现支持

获取服务器IP可使用ipv4包自动获取 https://github.com/go-basic/ipv4

用法：
```
go get github.com/go-basic/ipv4
ip = ipv4.LocalIP()
```
# 服务注册demo
```
package main

import (
	"github.com/go-basic/uuid"
	"github.com/go-zk/registry"
	"os"
	"os/signal"
	"syscall"
)

//服务注册示例
func main() {
	zkManager, err := registry.NewZkRegistry(
		registry.Hosts([]string{"127.0.0.1:2181"}),
		registry.Prefix("/zk-registry"),
		registry.Timeout(15),
		registry.Listens([]string{"test-service"}),
	)
	if err != nil {
		panic(err.Error())
	}
	defer zkManager.Close()
	node := &registry.Node{
		Id:      uuid.New(),
		Address: "127.0.0.1",
		Port:    8080,
	}
	zkManager.Register("test-service", node)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
```

# 服务发现demo

```
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
```
# http demo可自测更多功能
github.com/go-zk/registry/example/main.go