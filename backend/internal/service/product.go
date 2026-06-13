package service

import (
	"errors"

	"material-build/internal/model"
	"material-build/internal/repository"
)

type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductInput struct {
	MerchantID  uint64
	CategoryID  uint64
	Name        string
	Description string
	CoverImage  string
	Price       float64
	SalePrice   float64
	Unit        string
	Stock       int
}

func (s *ProductService) Create(in CreateProductInput) (*model.Product, error) {
	if in.Name == "" {
		return nil, errors.New("商品名称不能为空")
	}
	p := &model.Product{
		MerchantID:  in.MerchantID,
		CategoryID:  in.CategoryID,
		Name:        in.Name,
		Description: in.Description,
		CoverImage:  in.CoverImage,
		Price:       in.Price,
		SalePrice:   in.SalePrice,
		Unit:        in.Unit,
		Stock:       in.Stock,
		Status:      1,
	}
	if p.SalePrice == 0 {
		p.SalePrice = p.Price
	}
	if err := s.repo.CreateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Search(q repository.ProductQuery) ([]model.Product, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	return s.repo.ListProducts(q)
}

func (s *ProductService) Detail(id uint64) (*model.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) Categories() ([]model.Category, error) {
	return s.repo.ListCategories()
}
