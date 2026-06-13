package repository

import "material-build/internal/model"

type ProductQuery struct {
	Keyword    string
	CategoryID uint64
	MerchantID uint64
	CityID     uint64
	Page       int
	PageSize   int
}

func (r *Repository) ListProducts(q ProductQuery) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("products.status = 1")
	if q.Keyword != "" {
		query = query.Where("products.name LIKE ?", "%"+q.Keyword+"%")
	}
	if q.CategoryID > 0 {
		query = query.Where("products.category_id = ?", q.CategoryID)
	}
	if q.MerchantID > 0 {
		query = query.Where("products.merchant_id = ?", q.MerchantID)
	}
	if q.CityID > 0 {
		query = query.Joins("JOIN merchants ON merchants.id = products.merchant_id").
			Where("merchants.city_id = ?", q.CityID)
	}

	query.Count(&total)
	err := query.Preload("Merchant").Preload("Media").
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).
		Order("products.sales_count DESC").Find(&list).Error
	return list, total, err
}

func (r *Repository) GetProductByID(id uint64) (*model.Product, error) {
	var p model.Product
	err := r.db.Preload("Merchant").Preload("Media").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) CreateProduct(p *model.Product) error {
	return r.db.Create(p).Error
}

func (r *Repository) UpdateProduct(p *model.Product) error {
	return r.db.Save(p).Error
}

func (r *Repository) ListCategories() ([]model.Category, error) {
	var list []model.Category
	err := r.db.Where("status = 1").Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *Repository) AddProductMedia(media *model.ProductMedia) error {
	return r.db.Create(media).Error
}
