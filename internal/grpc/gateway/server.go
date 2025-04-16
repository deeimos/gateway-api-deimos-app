package gateway

import (
	"gateway-api/internal/lib/validation"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	"google.golang.org/grpc"
)

type serverApi struct {
	gateway_apiv1.UnimplementedPublicAuthServer
	gateway_apiv1.UnimplementedPublicServersServer
	gateway_apiv1.UnimplementedPublicMonitoringServer
	errMapper *validation.ErrorMapper
}

func Register(grpcServer *grpc.Server, errMapper *validation.ErrorMapper) {
	api := &serverApi{errMapper: errMapper}

	gateway_apiv1.RegisterPublicAuthServer(grpcServer, api)
	gateway_apiv1.RegisterPublicServersServer(grpcServer, api)
	gateway_apiv1.RegisterPublicMonitoringServer(grpcServer, api)
}
