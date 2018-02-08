package command

import (
	"context"

	grpc "google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/proxy"
)

type InboundOperation interface {
	ApplyInbound(context.Context, core.InboundHandler) error
}

type OutboundOperation interface {
	ApplyOutbound(context.Context, core.OutboundHandler) error
}

func (op *AddUserOperation) ApplyInbound(ctx context.Context, handler core.InboundHandler) error {
	getInbound, ok := handler.(proxy.GetInbound)
	if !ok {
		return newError("can't get inbound proxy from handler")
	}
	p := getInbound.GetInbound()
	um, ok := p.(proxy.UserManager)
	if !ok {
		return newError("proxy is not an UserManager")
	}
	return um.AddUser(ctx, op.User)
}

func (op *AddUserOperation) ApplyOutbound(ctx context.Context, handler core.OutboundHandler) error {
	getOutbound, ok := handler.(proxy.GetOutbound)
	if !ok {
		return newError("can't get outbound proxy from handler")
	}
	p := getOutbound.GetOutbound()
	um, ok := p.(proxy.UserManager)
	if !ok {
		return newError("proxy in not an UserManager")
	}
	return um.AddUser(ctx, op.User)
}

type handlerServer struct {
	s   *core.Instance
	ihm core.InboundHandlerManager
	ohm core.OutboundHandlerManager
}

func (s *handlerServer) AddInbound(ctx context.Context, request *AddInboundRequest) (*AddInboundResponse, error) {
	rawHandler, err := s.s.CreateObject(request.Inbound)
	if err != nil {
		return nil, err
	}
	handler, ok := rawHandler.(core.InboundHandler)
	if !ok {
		return nil, newError("not an InboundHandler.")
	}
	return &AddInboundResponse{}, s.ihm.AddHandler(ctx, handler)
}

func (s *handlerServer) RemoveInbound(ctx context.Context, request *RemoveInboundRequest) (*RemoveInboundResponse, error) {
	return &RemoveInboundResponse{}, s.ihm.RemoveHandler(ctx, request.Tag)
}

func (s *handlerServer) AlterInbound(ctx context.Context, request *AlterInboundRequest) (*AlterInboundResponse, error) {
	rawOperation, err := request.Operation.GetInstance()
	if err != nil {
		return nil, newError("unknown operation").Base(err)
	}
	operation, ok := rawOperation.(InboundOperation)
	if !ok {
		return nil, newError("not an inbound operation")
	}

	handler, err := s.ihm.GetHandler(ctx, request.Tag)
	if err != nil {
		return nil, newError("failed to get handler: ", request.Tag).Base(err)
	}

	return &AlterInboundResponse{}, operation.ApplyInbound(ctx, handler)
}

func (s *handlerServer) AddOutbound(ctx context.Context, request *AddOutboundRequest) (*AddOutboundResponse, error) {
	rawHandler, err := s.s.CreateObject(request.Outbound)
	if err != nil {
		return nil, err
	}
	handler, ok := rawHandler.(core.OutboundHandler)
	if !ok {
		return nil, newError("not an OutboundHandler.")
	}
	return &AddOutboundResponse{}, s.ohm.AddHandler(ctx, handler)
}

func (s *handlerServer) RemoveOutbound(ctx context.Context, request *RemoveOutboundRequest) (*RemoveOutboundResponse, error) {
	return &RemoveOutboundResponse{}, s.ohm.RemoveHandler(ctx, request.Tag)
}

func (s *handlerServer) AlterOutbound(ctx context.Context, request *AlterOutboundRequest) (*AlterOutboundResponse, error) {
	rawOperation, err := request.Operation.GetInstance()
	if err != nil {
		return nil, newError("unknown operation").Base(err)
	}
	operation, ok := rawOperation.(OutboundOperation)
	if !ok {
		return nil, newError("not an outbound operation")
	}

	handler := s.ohm.GetHandler(request.Tag)
	return &AlterOutboundResponse{}, operation.ApplyOutbound(ctx, handler)
}

type feature struct{}

func (*feature) Start() error {
	return nil
}

func (*feature) Close() error {
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := core.FromContext(ctx)
		if s == nil {
			return nil, newError("V is not in context.")
		}
		s.Commander().RegisterService(func(server *grpc.Server) {
			RegisterHandlerServiceServer(server, &handlerServer{
				s:   s,
				ihm: s.InboundHandlerManager(),
				ohm: s.OutboundHandlerManager(),
			})
		})
		return &feature{}, nil
	}))
}
