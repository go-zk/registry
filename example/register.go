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
	node := &registry.Node{
		Id:      uuid.New(),
		Address: "127.0.0.1",
		Port:    8080,
	}
	zkRegistry.Register("test-service", node)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
