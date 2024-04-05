package gatewayservice

import (
	"context"
	"fmt"

	omspb "github.com/ilivestrong/oms-gateway/internal/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	GatewayService struct {
		omspb.UnimplementedGatewayServiceServer

		productSvc omspb.ProductServiceClient
		orderSvc   omspb.OrderServiceClient
	}
)

func New(productSvcClient omspb.ProductServiceClient, orderSvcClient omspb.OrderServiceClient) *GatewayService {
	return &GatewayService{
		productSvc: productSvcClient,
		orderSvc:   orderSvcClient,
	}
}

func (gw *GatewayService) ListOrders(ctx context.Context, in *emptypb.Empty) (*omspb.ListOrdersResponse, error) {
	return gw.orderSvc.List(ctx, nil)
}

func (gw *GatewayService) CreateOrder(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	return gw.orderSvc.Create(ctx, req)
}

func (gw *GatewayService) GetProduct(ctx context.Context, req *omspb.GetProductRequest) (*omspb.Product, error) {
	return gw.productSvc.Get(ctx, req)
}

func (gw *GatewayService) ListProducts(ctx context.Context, req *omspb.ListProductsRequest) (*omspb.ListProductsResponse, error) {
	fmt.Println("calling ListProducts....")
	return gw.productSvc.List(ctx, req)
}

func (gw *GatewayService) CreateProduct(ctx context.Context, req *omspb.CreateProductRequest) (*omspb.Product, error) {
	return gw.productSvc.Create(ctx, req)
}

func (gw *GatewayService) UpdateProduct(ctx context.Context, req *omspb.UpdateProductRequest) (*omspb.Product, error) {
	return gw.productSvc.Update(ctx, req)
}

func (gw *GatewayService) DeleteProduct(ctx context.Context, req *omspb.DeleteProductRequest) (*omspb.DeleteProductResponse, error) {
	return gw.productSvc.Delete(ctx, req)
}

func (gw *GatewayService) DecrementProductQty(ctx context.Context, req *omspb.DecrementQtyRequest) (*omspb.DecrementQtyResponse, error) {
	return gw.productSvc.DecrementQty(ctx, req)
}
