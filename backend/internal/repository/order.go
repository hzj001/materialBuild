package repository

import "material-build/internal/model"

func (r *Repository) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *Repository) CreateOrderItems(items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *Repository) CreateDelivery(d *model.Delivery) error {
	return r.db.Create(d).Error
}

func (r *Repository) ListOrdersByUser(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64
	q := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	q.Count(&total)
	err := q.Preload("Items").Preload("Delivery").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *Repository) ListOrdersByMerchant(merchantID uint64, page, pageSize int) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64
	q := r.db.Model(&model.Order{}).Where("merchant_id = ?", merchantID)
	q.Count(&total)
	err := q.Preload("Items").Preload("Delivery").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *Repository) GetOrderByID(id uint64) (*model.Order, error) {
	var o model.Order
	err := r.db.Preload("Items").Preload("Delivery").First(&o, id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) UpdateOrderStatus(id uint64, status int8) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) CreateAfterSale(a *model.AfterSale) error {
	return r.db.Create(a).Error
}

func (r *Repository) ListAfterSalesByMerchant(merchantID uint64, page, pageSize int) ([]model.AfterSale, int64, error) {
	var list []model.AfterSale
	var total int64
	q := r.db.Model(&model.AfterSale{}).Where("merchant_id = ?", merchantID)
	q.Count(&total)
	err := q.Offset((page-1)*pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
