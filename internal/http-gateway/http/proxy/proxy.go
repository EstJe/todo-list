package proxy

import (
	"context"
	todoapi "github.com/EstJe/todo-list/api/gen"
	"github.com/EstJe/todo-list/internal/lib/op"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type Proxy struct {
	log      *slog.Logger
	mux      *runtime.ServeMux
	grpcAddr string
}

func New(log *slog.Logger, mux *runtime.ServeMux, grpcAddr string) *Proxy {
	return &Proxy{log: log, mux: mux, grpcAddr: grpcAddr}
}

func (p *Proxy) RegisterHandlers() {
	oper := op.Operation()

	log := p.log.With(
		slog.String("op", oper),
	)

	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	log.Info("registering handler", "GRPC address", p.grpcAddr)

	err := todoapi.RegisterTodoServiceHandlerFromEndpoint(ctx, p.mux, p.grpcAddr, opts)
	if err != nil {
		log.Error("failed to register handler", "error", err.Error())
		panic(op.Wrap(err))
	}

	log.Info("registered handler", "GRPC address", p.grpcAddr)
}
