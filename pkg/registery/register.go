package registery

import (
	"context"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcd "go.etcd.io/etcd/client/v3"
	ioetcd "ims-server/pkg/etcd"
	iologger "ims-server/pkg/logger"
	"strings"
	"sync"
)

const (
	ServiceRootPath = "/ims/http" // prefix for etcd key
	UserService     = "user"
)

// The following global variables are not visible outside the package and can only be accessed through GetServiceHub()
var (
	serviceHub *ServiceHub
	hubOnce    sync.Once // singleton instance, ensures that this global variable is only initialized once
)

type ServiceHub struct {
	client             *etcd.Client
	heartbeatFrequency int64 // heartbeat frequency (interval between heartbeats)
}

// Initialize ServiceHub using singleton pattern
func GetServiceHub(heartbeatFrequency int64) *ServiceHub {
	hubOnce.Do(func() {
		if serviceHub == nil {
			client, err := ioetcd.NewClient()
			if err != nil {
				iologger.Fatalf("Failed to connect to etcd: %v", err)
			}
			serviceHub = &ServiceHub{
				client:             client,
				heartbeatFrequency: heartbeatFrequency,
			}
		}
	})
	return serviceHub
}

// Register a service. For the first registration, write a key to etcd. For subsequent registrations, only renew the lease.
func (hub *ServiceHub) RegisterService(service string, endpoint string, leaseID etcd.LeaseID) (etcd.LeaseID, error) {
	ctx := context.Background()
	// If the lease ID is less than or equal to 0, the service is not registered yet, so automatically create a lease.
	if leaseID <= 0 {
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			return 0, err
		} else {
			key := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/" + endpoint // TrimRight is used to remove "/", ensuring that the subsequent addition is not duplicated
			if _, err := hub.client.Put(ctx, key, endpoint, etcd.WithLease(lease.ID)); err != nil {
				if err == rpctypes.ErrLeaseNotFound {
					return hub.RegisterService(service, endpoint, 0) // Fake lease, go through the registration process (set leaseID to 0)
				}
				return 0, err
			}
			return lease.ID, nil
		}
	}
	// If the lease already exists, simply renew it.
	if _, err := hub.client.KeepAliveOnce(ctx, leaseID); err != nil {
		iologger.Warn("Lease renewal failed: %v", err)
		return 0, err
	}
	return leaseID, nil
}

// Unregister a service.
func (hub *ServiceHub) UnRegisterService(service string, endpoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/" + endpoint

	if _, err := hub.client.Delete(ctx, key); err != nil {
		iologger.Warn("Failed to unregister service %s at endpoint %s: %v", service, endpoint, err)
		return err
	}
	iologger.Info("Unregistered service %s at endpoint %s", service, endpoint)
	return nil
}
