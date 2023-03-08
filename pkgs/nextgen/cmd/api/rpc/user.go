package rpc

import (
	"context"
	"fish_net/kitex_gen/user"
	"fish_net/kitex_gen/user/userservice"
	"fish_net/pkg/consts"
	"fish_net/pkg/errno"
	"fish_net/pkg/mw"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

var userClient userservice.Client

func initUser() {
	// r, err := etcd.NewEtcdResolver([]string{consts.ETCDAddress})
	// if err != nil {
	// 	panic(err)
	// }
	// provider.NewOpenTelemetryProvider(
	// 	provider.WithServiceName(consts.ApiServiceName),
	// 	provider.WithExportEndpoint(consts.ExportEndpoint),
	// 	provider.WithInsecure(),
	// )
	c, err := userservice.NewClient(
		consts.UserServiceName,
		// client.WithResolver(r),
		client.WithMuxConnection(1),
		client.WithMiddleware(mw.CommonMiddleware),
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.ApiServiceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

// RegisterBegin create user info
func RegisterBegin(ctx context.Context, req *user.CreateUserRequest) error {
	resp, err := userClient.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}
