package ioetcd

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
	"time"
)

func NewClient() (*etcd.Client, error) {
	// If etcd is not connected for 5 seconds, it will time out.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// connect to etcd
	var client *etcd.Client
	var err error
	go func() {
		config := ioconfig.GetEtcdConf()
		client, err = etcd.New(
			etcd.Config{
				Endpoints:   config.Endpoints,
				DialTimeout: config.DialTimeout,
			},
		)
		if err != nil {
			iologger.Fatalf("etcd client connect failed, err: %v", err)
		}
	}()
	// Check if the connection has timed out
	select {
	case <-ctx.Done():
		iologger.Fatalf("etcd client connect timeout")
	}
	return client, nil
}
