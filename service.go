package registry

import "time"

//服务
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

//节点
type Node struct {
	Id      string `json:"id"`      //uuid
	Name    string `json:"name"`    //服务名称
	Address string `json:"address"` //服务IP
	Port    int    `json:"port"`    //端口
	//下面是java使用的字段
	Payload            *Payload `json:"payload"`
	RegstrationTimeUTC int64    `json:"registrationTimeUTC"`
	ServiceType        string   `json:"serviceType"`
	UriSpec            *UriSpec `json:"uriSpec"`
}

/*********** 以下 java 使用 **************/
type Payload struct {
	Class string `json:"@class"`
}
type UriSpec struct {
	Parts []*UriSpecParts `json:"parts"`
}
type UriSpecParts struct {
	Value    string `json:"value"`
	Variable bool   `json:"variable"`
}

//添加java支持
func addJava(node *Node) {
	node.Payload = &Payload{
		Class: "org.springframework.cloud.zookeeper.discovery.ZookeeperInstance",
	}
	node.RegstrationTimeUTC = time.Now().UnixNano() / 1e6
	node.ServiceType = "DYNAMIC"
	node.UriSpec = &UriSpec{
		Parts: []*UriSpecParts{
			&UriSpecParts{
				Value:    "scheme",
				Variable: true,
			},
			&UriSpecParts{
				Value:    "://",
				Variable: false,
			},
			&UriSpecParts{
				Value:    "address",
				Variable: true,
			},
			&UriSpecParts{
				Value:    ":",
				Variable: true,
			},
			&UriSpecParts{
				Value:    "port",
				Variable: true,
			},
		},
	}

}
