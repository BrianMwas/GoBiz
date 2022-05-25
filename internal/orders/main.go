package orders

import (
	"awesomeProject/internal/entity"
	pb2 "awesomeProject/internal/orders/pb"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/utils"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Ptr Awesome generics in golang
func Ptr[T any](v T) *T {
	return &v
}

type OrderServer struct {
	pb2.UnimplementedOrdersServer
	Log   *logrus.Logger
	Store *db.MongoStore
}

func (s *OrderServer) GetOrders(ctx context.Context, req *pb2.GetOrdersReq) (*pb2.GetOrdersRes, error) {
	customerId := req.GetCustomerId()
	if len(customerId) > 0 {
		fmt.Println("We have a customer id")
		var orders []*entity.Order
		parseId, _ := primitive.ObjectIDFromHex(req.GetCustomerId())
		err := s.Store.GetAll(entity.OrderCollectionName, utils.KeyValue{"customerId": parseId}, &orders)
		if err != nil {
			logrus.Println("Failed to get all orders", err)
		}
		var pbOrders = make([]*pb2.Order, 0)

		for _, order := range orders {
			var products []*pb2.Product

			for _, item := range order.Items {
				products = append(products, &pb2.Product{
					Id:          Ptr(item.ID.Hex()),
					Price:       item.Price,
					Name:        item.Name,
					Description: &item.Description,
					Sku:         item.SKU,
				})
			}

			pbOrders = append(pbOrders, &pb2.Order{
				Id:           Ptr(order.ID.Hex()),
				Status:       order.Status,
				Items:        products,
				DeliveryDate: Ptr(int64(order.DeliveryDate.T)),
				OrderNo:      order.OrderNo,
				CustomerId:   order.CustomerId.Hex(),
			})
		}

		return &pb2.GetOrdersRes{Orders: pbOrders}, nil
	} else {
		var orders []entity.Order
		err := s.Store.GetAll(entity.OrderCollectionName, utils.KeyValue{}, &orders)
		fmt.Println("Orders slice ", orders)
		if err != nil {
			return nil, err
		}
		var pbOrders = make([]*pb2.Order, 0)

		for _, order := range orders {
			var products = make([]*pb2.Product, 0)

			for _, item := range order.Items {
				products = append(products, &pb2.Product{
					Id:          Ptr(item.ID.Hex()),
					Price:       item.Price,
					Name:        item.Name,
					Description: &item.Description,
					Sku:         item.SKU,
				})
			}

			pbOrders = append(pbOrders, &pb2.Order{
				Id:           Ptr(order.ID.Hex()),
				Status:       order.Status,
				Items:        products,
				DeliveryDate: Ptr(int64(order.DeliveryDate.T)),
				OrderNo:      order.OrderNo,
				CustomerId:   order.CustomerId.Hex(),
			})
		}

		return &pb2.GetOrdersRes{Orders: pbOrders}, nil
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb2.CreateOrderReq) (*pb2.CreateOrderRes, error) {
	newOrder := entity.NewOrder()
	reqOrder := req.GetOrder()

	newOrder.OrderNo = reqOrder.GetOrderNo()
	newOrder.Status = reqOrder.GetStatus()
	newOrder.CustomerId, _ = primitive.ObjectIDFromHex(reqOrder.GetCustomerId())
	newOrder.DeliveryDate = primitive.Timestamp{
		T: uint32(time.UnixMilli(reqOrder.GetDeliveryDate()).Unix()),
		I: 0,
	}
	var products = make([]*entity.Product, 0)
	for _, item := range reqOrder.Items {
		newProduct := entity.NewProduct()
		newProduct.ID = primitive.NewObjectID()
		newProduct.Price = item.GetPrice()
		newProduct.Name = item.GetName()
		newProduct.SKU = item.GetSku()
		newProduct.Description = item.GetDescription()
		newProduct.CreatedAt = primitive.Timestamp{
			T: uint32(time.Now().Unix()),
			I: 0,
		}
		products = append(products, newProduct)
	}

	newOrder.Items = products
	order, err := newOrder.Persist(s.Store)

	if err != nil {
		return nil, err
	}

	if err != nil {
		fmt.Println("Error creating a new Order", err)
		return nil, err
	}
	var items []*pb2.Product
	for _, item := range newOrder.Items {
		items = append(items, &pb2.Product{
			Id:          Ptr(item.ID.Hex()),
			Price:       item.Price,
			Name:        item.Name,
			Description: &item.Description,
		})
	}

	return &pb2.CreateOrderRes{Order: &pb2.Order{
		Id:           Ptr(order.ID.Hex()),
		Status:       newOrder.Status,
		Items:        items,
		DeliveryDate: Ptr(int64(newOrder.DeliveryDate.T)),
	}}, nil
}

func (s *OrderServer) CreateCustomer(ctx context.Context, req *pb2.CreateCustomerReq) (*pb2.CreateCustomerRes, error) {
	s.Log.Println("Creating a new customer ", req.GetLat())
	validate := validator.New()
	newCustomer := entity.NewCustomer()
	fmt.Println("Customer ", req.GetName())

	newCustomer.Name = req.GetName()
	p := entity.Location{
		Lat:  req.GetLat(),
		Long: req.GetLon(),
	}

	newCustomer.Position = p
	err := validate.Struct(newCustomer)

	if err != nil {
		return &pb2.CreateCustomerRes{
			Success: false,
			Message: err.Error(),
		}, err
	}

	_, insertErr := s.Store.Insert(entity.CustomerCollectionName, newCustomer)

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
	orderId := req.GetId()

	order, err := entity.NewOrder().Get(s.Store, utils.KeyValue{"_id": orderId})

	if err != nil {
		logrus.Println("We failed getting order", err)
	}
	order.Status = req.GetStatus()
	_, persistErr := order.Persist(s.Store)

	if persistErr != nil {
		return nil, persistErr
	}

	return &pb2.UpdateOrderStatusRes{
		Success: true,
		Message: "Update Successful",
	}, nil
}

func (s *OrderServer) UpdateOrder(ctx context.Context, req *pb2.UpdateOrderReq) (*pb2.UpdateOrderRes, error) {
	updateOrder := req.GetOrder()
	order, err := entity.NewOrder().Get(s.Store, utils.KeyValue{
		"_id": updateOrder.GetId(),
	})

	if err != nil {
		return nil, err
	}
	var products = make([]*entity.Product, 0)
	for _, item := range updateOrder.Items {
		newProduct := entity.NewProduct()

		newProduct.Price = item.GetPrice()
		newProduct.Name = item.GetName()
		newProduct.Description = item.GetDescription()
		newProduct.Name = item.GetSku()
		products = append(products, newProduct)
	}

	order.Items = products

	_, persistErr := order.Persist(s.Store)

	if persistErr != nil {
		return nil, persistErr
	}

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
