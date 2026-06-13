package repository

import "material-build/internal/model"

func (r *Repository) FindMerchantByUserID(userID uint64) (*model.Merchant, error) {
	var m model.Merchant
	err := r.db.Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) ListMerchants(cityID uint64, page, pageSize int) ([]model.Merchant, int64, error) {
	var list []model.Merchant
	var total int64
	q := r.db.Model(&model.Merchant{}).Where("status = 1")
	if cityID > 0 {
		q = q.Where("city_id = ?", cityID)
	}
	q.Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("rating DESC").Find(&list).Error
	return list, total, err
}

func (r *Repository) CreateMerchant(m *model.Merchant) error {
	return r.db.Create(m).Error
}

func (r *Repository) UpdateMerchantStatus(id uint64, status int8) error {
	return r.db.Model(&model.Merchant{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) ListPendingMerchants(page, pageSize int) ([]model.Merchant, int64, error) {
	var list []model.Merchant
	var total int64
	q := r.db.Model(&model.Merchant{}).Where("status = 0")
	q.Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
