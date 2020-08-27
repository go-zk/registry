package main

import (
	"fmt"
	"github.com/go-basic/uuid"
	"github.com/go-zk/registry"
	"log"
	"math/rand"
	"net/http"
)

var (
	zkRegistry *registry.ZkRegistry
	err        error
	port       = "8080"
)

//注册中心http demo
func main() {
	zkRegistry, err = registry.NewZkRegistry(
		registry.Hosts([]string{"127.0.0.1:2181"}),
		registry.Prefix("/zk-registry"),
		registry.Timeout(5),
		registry.Listens([]string{"test-service"}),
	)
	if err != nil {
		panic(err.Error())
	}
	defer zkRegistry.Close()
	//随机获取一个服务节点
	http.HandleFunc("/", getServiceNode)
	//注册一个节点到test-service
	http.HandleFunc("/register", register)
	//反注册一个节点test-service
	http.HandleFunc("/deregister", deregister)
	//从test-service删除指定节点（不会从注册中心删除）
	http.HandleFunc("/del", del)
	//获取一个服务下的所有节点
	http.HandleFunc("/nodes", nodes)
	//获取监听服务列表
	http.HandleFunc("/service", service)
	log.Println("start running : " + port)
	http.ListenAndServe(port, nil)
}

//随机获取一个服务节点
func getServiceNode(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	svc := get.Get("service")
	if svc == "" {
		writer.Write([]byte("请输入服务名称"))
		return
	}
	node, err := zkRegistry.GetServerNode(svc)
	if err != nil {
		writer.Write([]byte(err.Error()))
	} else {
		writer.Write([]byte(fmt.Sprintf("%+v", node)))
	}
}

//注册test-service节点到注册中心
// ip 一个IP地址，端口随机
func register(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	ip := get.Get("ip")
	node := &registry.Node{
		Id:      uuid.New(),
		Address: ip,
		Port:    rand.Intn(8080),
	}
	zkRegistry.Register("test-service", node)
	if err != nil {
		writer.Write([]byte(err.Error()))
	} else {
		writer.Write([]byte(fmt.Sprintf("%+v", node)))
	}
}

//注册中心删除test-service的一个节点
// id 为nodeId
func deregister(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	id := get.Get("id")
	node := &registry.Node{
		Id:   id,
		Port: rand.Intn(8080),
	}

	zkRegistry.Deregister("test-service", node)
	if err != nil {
		writer.Write([]byte(err.Error()))
	} else {
		writer.Write([]byte(fmt.Sprintf("%+v", node)))
	}
}

//从test-service删除指定节点（不会从注册中心删除）
// id 为nodeId
func del(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	id := get.Get("id")
	node := &registry.Node{
		Id: id,
	}
	err := zkRegistry.DelServiceNode("test-service", node)
	if err != nil {
		writer.Write([]byte(err.Error()))
	} else {
		writer.Write([]byte(fmt.Sprintf("%+v", node)))
	}
}

//获取一个服务下的所有节点
func nodes(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	svc := get.Get("service")
	if svc == "" {
		svc = "test-service"
	}
	all := zkRegistry.GetAllNode(svc)
	for _, item := range all {
		writer.Write([]byte(fmt.Sprintf("%+v", item) + "\n"))
	}
}

//获取监听服务列表
func service(writer http.ResponseWriter, request *http.Request) {
	all := zkRegistry.GetAllService()
	for _, item := range all {
		writer.Write([]byte(fmt.Sprintf("%+v", item) + "\n"))
	}
}
