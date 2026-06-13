package repository

import "material-build/internal/model"

func (r *Repository) ListCities() ([]model.City, error) {
	var cities []model.City
	err := r.db.Where("status = 1").Order("id ASC").Find(&cities).Error
	return cities, err
}

func (r *Repository) CreateCity(city *model.City) error {
	return r.db.Create(city).Error
}

func (r *Repository) UpdateCityStatus(id uint64, status int8) error {
	return r.db.Model(&model.City{}).Where("id = ?", id).Update("status", status).Error
}
