package options

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Addrs       []string //required
	NamespaceId string   //required
	Timeout     time.Duration
	LogLevel    string
	RotateTime  time.Duration
}

type Option func(*Options)

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func NameSpaceId(namespaceId string) Option {
	return func(o *Options) {
		o.NamespaceId = namespaceId
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func RotateTime(t time.Duration) Option {
	return func(o *Options) {
		o.RotateTime = t
	}
}

func LogLevel(b string) Option {
	return func(o *Options) {
		o.LogLevel = b
	}
}

func OptionInit(opts ...Option) vo.NacosClientParam {
	opt := parseOptions(opts...)
	if len(opt.Addrs) == 0 {
		panic("请设置注册中心地址")
	}
	if opt.NamespaceId == "" {
		panic("请设置注册中心NamespaceId")
	}
	if opt.LogLevel == "" {
		opt.LogLevel = "debug"
	}
	if opt.Timeout == 0 {
		opt.Timeout = 5 * time.Second
	}
	if opt.RotateTime == 0 {
		opt.RotateTime = 1 * time.Hour
	}
	serverConfigs := parseServerConfig(opt)
	clientConfig := parseClientConfig(opt)

	return vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	}
}

func parseClientConfig(opts *Options) *constant.ClientConfig {
	clientConfig := constant.ClientConfig{
		NamespaceId:         opts.NamespaceId,
		TimeoutMs:           uint64(opts.Timeout.Milliseconds()),
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          opts.RotateTime.String(),
		MaxAge:              3,
		LogLevel:            opts.LogLevel,
	}
	return &clientConfig
}
func parseServerConfig(opt *Options) []constant.ServerConfig {
	var srvConfig []constant.ServerConfig
	for _, address := range opt.Addrs {
		path := "/nacos123"
		addr, port, err := net.SplitHostPort(address)
		ports := strings.Split(port, "/")
		if len(ports) > 1 {
			port = ports[0]
			path = "/" + strings.Join(ports[1:], "/")
		}
		_port, _ := strconv.ParseUint(port, 10, 64)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			_port = 8848
			addr = address
		}

		config := constant.ServerConfig{
			IpAddr:      addr,
			ContextPath: path,
			Port:        _port,
			Scheme:      "http",
		}
		srvConfig = append(srvConfig, config)
	}
	return srvConfig
}
func parseOptions(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}
