package orders

import (
	"awesomeProject/internal/entity"
	pb2 "awesomeProject/internal/orders/pb"
	"awesomeProject/pkg/db"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type OrderServer struct {
	pb2.UnimplementedOrdersServer
	Log   *logrus.Logger
	Store *db.MongoStore
}

func (s *OrderServer) GetOrderDetails(ctx context.Context, req *pb2.GetOrderDetailsReq) (*pb2.GetOrderDetailsRes, error) {
	s.Log.Println("Get order details")
	return &pb2.GetOrderDetailsRes{Order: &pb2.Order{
		Id:           "8989",
		Status:       pb2.OrderStatus_processing,
		Items:        nil,
		DeliveryDate: time.Now().UnixMilli(),
		CreatedAt:    time.Now().UnixMilli(),
	}}, nil
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb2.CreateOrderReq) (*pb2.CreateOrderRes, error) {
	s.Log.Println("Creating a new order")

	return &pb2.CreateOrderRes{Order: &pb2.Order{
		Id:           "8989",
		Status:       pb2.OrderStatus_processing,
		Items:        nil,
		DeliveryDate: time.Now().UnixMilli(),
		CreatedAt:    time.Now().UnixMilli(),
	}}, nil
}

func (s *OrderServer) CreateCustomer(ctx context.Context, req *pb2.CreateCustomerReq) (*pb2.CreateCustomerRes, error) {
	s.Log.Println("Creating a new customer ", req.GetLat())

	newCustomer := entity.NewCustomer()

	newCustomer.Name = req.GetName()
	newCustomer.Position.Lat = req.GetLat()
	newCustomer.Position.Long = req.GetLon()

	//err := validate.Struct(newCustomer)

	//if err != nil {
	//	return &pb2.CreateCustomerRes{
	//		Success: false,
	//		Message: err,
	//	}, err
	//}

	insertErr := s.Store.Insert(entity.CustomerCollectionName, newCustomer)

	if insertErr != nil {
		s.Log.Println("Failed to create a new customer")
		return nil, insertErr
	}

	return &pb2.CreateCustomerRes{
		Success: true,
		Message: "Customer created Successfully ",
	}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb2.UpdateOrderStatusReq) (*pb2.UpdateOrderStatusRes, error) {
	s.Log.Println("Update order status", req.Status)
	return &pb2.UpdateOrderStatusRes{
		Success: true,
		Message: "Update Successful",
	}, nil
}

func (s *OrderServer) UpdateOrder(ctx context.Context, req *pb2.UpdateOrderReq) (*pb2.UpdateOrderRes, error) {
	s.Log.Println("Updating a full order ", req.Order)
	return &pb2.UpdateOrderRes{Order: req.Order}, nil
}
func (s *OrderServer) MustEmbedUnimplementedOrdersServer() {}

func NewOrderServer(log *logrus.Logger, mongo *db.MongoStore) *OrderServer {
	return &OrderServer{
		UnimplementedOrdersServer: pb2.UnimplementedOrdersServer{},
		Log:                       log,
		Store:                     mongo,
	}
}
