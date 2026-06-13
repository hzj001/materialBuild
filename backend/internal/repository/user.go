package repository

import "material-build/internal/model"

func (r *Repository) FindUserByPhoneAndRole(phone, role string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ? AND role = ?", phone, role).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
