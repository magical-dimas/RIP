package repository

import (
	"fmt"
	"rip_project/internal/app/ds"
)

func (r *Repository) GetModels() ([]ds.Model, error) {
	var models []ds.Model
	err := r.db.Where("is_deleted = ?", false).Find(&models).Error
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("массив моделей пуст")
	}
	return models, nil
}

func (r *Repository) GetModel(id int) (ds.Model, error) {
	var model ds.Model
	err := r.db.Where("model_id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return ds.Model{}, err
	}
	return model, nil
}

func (r *Repository) GetModelsByTitle(title string) ([]ds.Model, error) {
	var models []ds.Model
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
