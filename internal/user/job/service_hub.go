package job

import (
	"context"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcd "go.etcd.io/etcd/client/v3"
	ioetcd "ims-server/pkg/etcd"
	iologger "ims-server/pkg/logger"
	"log"
	"strings"
	"sync"
)

// 以下全局变量包外不可见，包外想使用时通过 GetServiceHub() 获得
var (
	serviceHub *ServiceHub
	hubOnce    sync.Once // 单例实例，单例模式中确保该全局变量只被初始化一次
)

type ServiceHub struct {
	client             *etcd.Client
	heartbeatFrequency int64 // 心跳频率（心跳间隔时间）
}

// 初始化 ServiceHub，单例模式
func GetServiceHub(heartbeatFrequency int64) *ServiceHub {
	hubOnce.Do(func() {
		if serviceHub == nil {
			client, err := ioetcd.NewClient()
			if err != nil {
				iologger.Fatalf("连接etcd失败: %v", err)
			}
			serviceHub = &ServiceHub{
				client:             client,
				heartbeatFrequency: heartbeatFrequency,
			}
		}
	})
	return serviceHub
}

// 服务注册：第一次注册向 etcd 写一个 key，后续注册仅仅是在续约
// 参数：
// - service: 服务名称
// - endpoint: 服务地址
func (hub *ServiceHub) RegisterService(service string, endpoint string, leaseID etcd.LeaseID) (etcd.LeaseID, error) {
	ctx := context.Background()
	// 租约 ID 小于 0，服务未注册，自动创建一个租约
	if leaseID <= 0 {
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			return 0, err
		} else {
			key := strings.TrimRight(ioetcd.ServiceRootPath, "/") + "/" + service + "/" + endpoint // TrimRight 用于删除/，以确保后面的添加不会重复
			if _, err := hub.client.Put(ctx, key, endpoint, etcd.WithLease(lease.ID)); err != nil {
				if err == rpctypes.ErrLeaseNotFound {
					return hub.RegisterService(service, endpoint, 0) // 虚假租约，走注册流程 (把 leaseID 置为 0)
				}
				return 0, err
			}
			return lease.ID, nil
		}
	}
	// 已经存在租约，直接续租
	if _, err := hub.client.KeepAliveOnce(ctx, leaseID); err != nil {
		log.Printf("续约失败:%v", err)
		return 0, err
	}
	return leaseID, nil
}

// 注销服务
func (hub *ServiceHub) UnRegisterService(service string, endpoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(ioetcd.ServiceRootPath, "/") + "/" + service + "/" + endpoint

	if _, err := hub.client.Delete(ctx, key); err != nil {
		log.Printf("注销服务 %s 对应的节点 %s 失败: %v", service, endpoint, err)
		return err
	}
	log.Printf("注销服务 %s 对应的节点 %s", service, endpoint)
	return nil
}
