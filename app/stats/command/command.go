package command

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg command -path App,Stats,Command

import (
	"context"

	grpc "google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/common"
)

type statsServer struct {
	stats core.StatManager
}

func (s *statsServer) GetStats(ctx context.Context, request *GetStatsRequest) (*GetStatsResponse, error) {
	c := s.stats.GetCounter(request.Name)
	if c == nil {
		return nil, newError(request.Name, " not found.")
	}
	var value int64
	if request.Reset_ {
		value = c.Set(0)
	} else {
		value = c.Value()
	}
	return &GetStatsResponse{
		Stat: &Stat{
			Name:  request.Name,
			Value: value,
		},
	}, nil
}

type service struct {
	v *core.Instance
}

func (s *service) Register(server *grpc.Server) {
	RegisterStatsServiceServer(server, &statsServer{
		stats: s.v.Stats(),
	})
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := core.MustFromContext(ctx)
		return &service{v: s}, nil
	}))
}
