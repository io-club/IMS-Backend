package registery

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"ims-server/pkg/etcd"
	"log"
	"strings"
	"sync"
)

type ServiceClient struct {
	client        *etcd.Client
	endpointCache sync.Map // 维护每一个 service 下的所有 servers
	watched       sync.Map
}

var (
	serviceClient *ServiceClient
	clientOnce    sync.Once
)

// ServiceClient 的构造函数，单例模式
func GetServiceClient() *ServiceClient {
	clientOnce.Do(func() {
		if serviceClient == nil {
			if client, err := ioetcd.NewClient(); err != nil {
				log.Fatalf("连接不上etcd服务器: %v", err) //发生 log.Fatal 时 go 进程会直接退出
			} else {
				serviceClient = &ServiceClient{
					client:        client,
					endpointCache: sync.Map{},
					watched:       sync.Map{},
				}
			}
		}
	})
	return serviceClient
}

// 从 etcd 里获取服务对应的节点（endpoint）
func (hub *ServiceClient) getServiceEndpoints(service string) []string {
	ctx := context.Background()
	prefix := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcd.WithPrefix()); err != nil { //按前缀获取 key-value
		log.Printf("获取服务%s的节点失败: %v", service, err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			path := strings.Split(string(kv.Key), "/") // 只需要 key（服务地址）
			endpoints = append(endpoints, path[len(path)-1])
		}
		log.Printf("刷新%s服务对应的server %v\n", service, endpoints)
		return endpoints
	}
}

// 把第一次查询 etcd 的结果缓存起来，然后安装一个 Watcher，仅 etcd 数据变化时更新本地缓存，这样可以降低 etcd 的访问压力
func (hub *ServiceClient) watchEndpointsOfService(service string) {
	if _, ok := hub.watched.LoadOrStore(service, true); ok {
		return // 已经被监听过了
	}
	ctx := context.Background()
	prefix := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/"
	ch := hub.client.Watch(ctx, prefix, etcd.WithPrefix()) //根据前缀监听，每一个修改都会放入管道 ch
	log.Printf("监听服务%s的节点变化", service)
	go func() {
		for resp := range ch { // 与 _,response:=range ch 不同，这是一种特殊的 range 循环语法，会一直迭代管道直到关闭
			for _, event := range resp.Events { // 每次从 ch 里取出来的是事件的集合
				path := strings.Split(string(event.Kv.Key), "/") // 按 / 分割
				if len(path) > 2 {
					service := path[len(path)-2]
					endpoints := hub.getServiceEndpoints(service)
					if len(endpoints) > 0 {
						hub.endpointCache.Store(service, endpoints) // 将 service 对应的 endpoints 存入 map 作为缓存
					} else {
						hub.endpointCache.Delete(service) // 该服务已没有可用的节点，删除键值对（缓存）
					}
				}
			}
		}
	}()
}

// 服务发现：client 每次进行 RPC 调用之前都查询 etcd，获取 server 集合，然后采用负载均衡算法选择一台 server
func (hub *ServiceClient) GetServiceEndpointsWithCache(service string) []string {
	hub.watchEndpointsOfService(service)
	if endpoints, ok := hub.endpointCache.Load(service); ok {
		return endpoints.([]string)
	}
	endpoints := hub.getServiceEndpoints(service)
	if len(endpoints) > 0 {
		hub.endpointCache.Store(service, endpoints) // 将 service 对应的 endpoints 存入 map 作为缓存
	}
	return endpoints
}
