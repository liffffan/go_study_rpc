package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go_study_rpc/registry"
	"log"
	"path"
	"time"
)

const MaxServiceNum = 8

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	// 定义一个 map 来存储 service
	registryServiceMap map[string]*RegisterService
}

type RegisterService struct {
	id      clientv3.LeaseID
	service *registry.Service
	// 用一个变量来标明有没有注册过
	registered bool
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh: make(chan *registry.Service, MaxServiceNum),
		// 初始化map
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

func init() {
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

// 插件的名字
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

// 初始化
func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {

	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})

	if err != nil {
		return
	}

	return
}

// 服务注册
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {

	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}

	return
}

func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (e *EtcdRegistry) servicePath(service registry.Service) string {
	return path.Join(e.options.RegistryPath, service.Name)
}

func (e *EtcdRegistry) run() {
	// 首先检测管道里有没有注册，注册的话就拿出来，按期续约就可以了，没有注册就拿出来注册
	// 从管道里把要注册的服务拿出来
	select {
	case service := <-e.serviceCh:
		// 需要判断 map 里是否已经存在这个服务了
		_, ok := e.registryServiceMap[service.Name]
		if ok {
			break
		}
		// 如果不存在就放入 map 里
		registryService := &RegisterService{
			service: service,
		}
		e.registryServiceMap[service.Name] = registryService
	default:
		// 做续约操作
		e.registerOrKeepAlive()
	}

	// 连接 etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 首先获取租约的授权
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// 然后通过租约的 id 来保持续期
	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 永久续约
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}

	for {
		ka := <-ch
		fmt.Println("ttl:", ka)
	}
}

func (e *EtcdRegistry) registerOrKeepAlive() {
	// 遍历 map 里的所有服务，进行续约

}
