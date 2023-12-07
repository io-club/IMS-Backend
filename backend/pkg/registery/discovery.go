package registery

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"ims-server/pkg/etcd"
	iologger "ims-server/pkg/logger"
	"strings"
	"sync"
)

type ServiceClient struct {
	client        *etcd.Client
	endpointCache sync.Map // maintains all servers under each service
	watched       sync.Map
}

var (
	serviceClient *ServiceClient
	clientOnce    sync.Once
)

// GetServiceClient is the constructor for ServiceClient using singleton pattern
func GetServiceClient() *ServiceClient {
	clientOnce.Do(func() {
		if serviceClient == nil {
			if client, err := ioetcd.NewClient(); err != nil {
				iologger.Fatalf("Failed to connect to etcd server: %v", err) // When log.Fatal is called, the Go process exits directly
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

// Get service endpoints from etcd
func (hub *ServiceClient) getServiceEndpoints(service string) []string {
	ctx := context.Background()
	prefix := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcd.WithPrefix()); err != nil {
		iologger.Warn("Failed to get endpoints for service %s: %v", service, err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			path := strings.Split(string(kv.Key), "/") // Only need the key (service address)
			endpoints = append(endpoints, path[len(path)-1])
		}
		iologger.Info("Refreshed servers for service %s: %v\n", service, endpoints)
		return endpoints
	}
}

// Cache the initial query result from etcd and then install a watcher to update the local cache only when etcd data changes, reducing the pressure on etcd access
func (hub *ServiceClient) watchEndpointsOfService(service string) {
	if _, ok := hub.watched.LoadOrStore(service, true); ok {
		return // Already being watched
	}
	ctx := context.Background()
	prefix := strings.TrimRight(ServiceRootPath, "/") + "/" + service + "/"
	ch := hub.client.Watch(ctx, prefix, etcd.WithPrefix()) // Watch based on prefix, each modification will be put into the channel ch
	iologger.Info("Watching for endpoint changes of service %s", service)
	go func() {
		for resp := range ch {
			for _, event := range resp.Events {
				path := strings.Split(string(event.Kv.Key), "/") // Split by /
				if len(path) > 2 {
					service := path[len(path)-2]
					endpoints := hub.getServiceEndpoints(service)
					if len(endpoints) > 0 {
						hub.endpointCache.Store(service, endpoints) // Store the endpoints for the service in the map as cache
					} else {
						hub.endpointCache.Delete(service) // No available endpoints for the service, delete the key-value pair (cache)
					}
				}
			}
		}
	}()
}

// GetServiceEndpointsWithCache queries etcd for the server collection before each RPC call, and then selects a server using a load balancing algorithm
func (hub *ServiceClient) GetServiceEndpointsWithCache(service string) []string {
	hub.watchEndpointsOfService(service)
	if endpoints, ok := hub.endpointCache.Load(service); ok {
		return endpoints.([]string)
	}
	endpoints := hub.getServiceEndpoints(service)
	if len(endpoints) > 0 {
		hub.endpointCache.Store(service, endpoints) // Cache endpoints
	}
	return endpoints
}
