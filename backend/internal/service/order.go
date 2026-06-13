package service

import (
	"errors"
	"fmt"
	"time"

	"material-build/internal/model"
	"material-build/internal/repository"
)

type OrderService struct {
	repo *repository.Repository
}

func NewOrderService(repo *repository.Repository) *OrderService {
	return &OrderService{repo: repo}
}

type OrderItemInput struct {
	ProductID uint64
	Quantity  int
}

type CreateOrderInput struct {
	UserID         uint64
	MerchantID     uint64
	Items          []OrderItemInput
	ReceiverName   string
	ReceiverPhone  string
	Address        string
	Remark         string
}

func (s *OrderService) Create(in CreateOrderInput) (*model.Order, error) {
	if len(in.Items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}

	var total float64
	var orderItems []model.OrderItem

	for _, item := range in.Items {
		product, err := s.repo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("商品不存在: %d", item.ProductID)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("商品 %s 库存不足", product.Name)
		}
		subtotal := product.SalePrice * float64(item.Quantity)
		total += subtotal
		orderItems = append(orderItems, model.OrderItem{
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductImage: product.CoverImage,
			Price:        product.SalePrice,
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
		})
	}

	orderNo := fmt.Sprintf("MB%s%04d", time.Now().Format("20060102150405"), in.UserID%10000)
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      in.UserID,
		MerchantID:  in.MerchantID,
		TotalAmount: total,
		PayAmount:   total,
		Status:      0,
		Remark:      in.Remark,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	if err := s.repo.CreateOrderItems(orderItems); err != nil {
		return nil, err
	}

	delivery := &model.Delivery{
		OrderID:       order.ID,
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Address:       in.Address,
		DeliveryType:  "merchant",
		Status:        0,
	}
	if err := s.repo.CreateDelivery(delivery); err != nil {
		return nil, err
	}

	order.Items = orderItems
	order.Delivery = delivery
	return order, nil
}

func (s *OrderService) ListByUser(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return s.repo.ListOrdersByUser(userID, page, pageSize)
}

func (s *OrderService) ListByMerchant(merchantID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return s.repo.ListOrdersByMerchant(merchantID, page, pageSize)
}

func (s *OrderService) Pay(orderID uint64) error {
	now := time.Now()
	return s.repo.DB().Model(&model.Order{}).Where("id = ? AND status = 0", orderID).
		Updates(map[string]interface{}{
			"status":     1,
			"pay_method": "wechat",
			"pay_at":     now,
		}).Error
}

func (s *OrderService) UpdateStatus(orderID uint64, status int8) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}
