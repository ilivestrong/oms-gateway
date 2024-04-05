package internal

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	OrderSvcClientConn   *grpc.ClientConn
	ProductSvcClientConn *grpc.ClientConn
}

func New(opts *Options) (*Service, error) {
	svc := &Service{}
	if err := initializeRpcConnections(opts, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func initializeRpcConnections(opts *Options, svc *Service) error {
	conn, err := grpc.Dial(opts.OrderServiceListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return err
	}
	svc.OrderSvcClientConn = conn

	conn, err = grpc.Dial(opts.ProductServiceListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//log.Fatalf("could not connect: %v", err)
		return err
	}
	svc.ProductSvcClientConn = conn
	return nil
}

func (svc *Service) Shutdown(ctx context.Context) error {
	return nil
}
