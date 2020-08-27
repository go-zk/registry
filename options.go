package registry

type Options struct {
	hosts   []string
	prefix  string
	timeout int
	//监听的service(需要在注册中心注册)
	listens []string
}

//默认值
var (
	DefaultHosts   = []string{"127.0.0.1:2181"}
	DefaultPrefix  = "/zk-registry"
	DefaultTimeout = 5
)

type Option func(*Options)

func (z *ZkRegistry) Options() Options {
	return z.opts
}

func newOptions(opts ...Option) Options {
	opt := Options{
		hosts:   DefaultHosts,
		prefix:  DefaultPrefix,
		timeout: DefaultTimeout,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

//zk主机列表
func Hosts(hosts []string) Option {
	return func(o *Options) {
		o.hosts = hosts
	}
}

//注册中心前缀
func Prefix(prefix string) Option {
	return func(o *Options) {
		o.prefix = prefix
	}
}

//链接超时时间
func Timeout(timeout int) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

//监听的服务列表
func Listens(listens []string) Option {
	return func(o *Options) {
		o.listens = listens
	}
}
