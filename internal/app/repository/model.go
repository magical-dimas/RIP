package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rip_project/internal/app/ds"
	minioClient "rip_project/internal/app/minioClient"
	"rip_project/internal/app/serializer"
)

func (r *Repository) GetModels() ([]ds.Model, error) {
	var models []ds.Model
	err := r.db.Where("is_deleted = ?", false).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *Repository) GetModel(id int) (*ds.Model, error) {
	var model ds.Model
	err := r.db.Where("model_id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: модель с id %d", ErrNotFound, id)
		}
		return nil, err
	}
	return &model, nil
}

func (r *Repository) GetModelsByTitle(title string) ([]ds.Model, error) {
	var models []ds.Model
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *Repository) CreateModel(j serializer.ModelJSON) (ds.Model, error) {
	model := serializer.ModelFromJSON(j)
	err := r.db.Create(&model).Scan(&model).Error
	if err != nil {
		return ds.Model{}, err
	}
	return model, nil
}

func (r *Repository) AddPhoto(ctx *gin.Context, modelID int, file *multipart.FileHeader) (*ds.Model, error) {
	model, err := r.GetModel(modelID)
	if err != nil {
		return nil, err
	}
	if model.PhotoURL != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), model.PhotoURL)
	}
	fileName, err := minioClient.UploadImage(ctx, r.mc, minioClient.GetImgBucket(), file, model.ModelID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&ds.Model{}).Where("model_id = ?", modelID).Update("photo_url", fileName).Error; err != nil {
		return nil, err
	}
	model.PhotoURL = fileName
	return model, nil
}

func (r *Repository) AddVideo(ctx *gin.Context, modelID int, file *multipart.FileHeader) (*ds.Model, error) {
	model, err := r.GetModel(modelID)
	if err != nil {
		return nil, err
	}
	if model.Video != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), model.Video)
	}
	fileName, err := minioClient.UploadVideo(ctx, r.mc, minioClient.GetImgBucket(), file, model.ModelID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&ds.Model{}).Where("model_id = ?", modelID).Update("video", fileName).Error; err != nil {
		return nil, err
	}
	model.Video = fileName
	return model, nil
}

func (r *Repository) DeleteModel(id int) error {
	model, err := r.GetModel(id)
	if err != nil {
		return err
	}
	if model.PhotoURL != "" {
		_ = minioClient.DeleteObject(context.Background(), r.mc, minioClient.GetImgBucket(), model.PhotoURL)
	}
	if model.Video != "" {
		_ = minioClient.DeleteObject(context.Background(), r.mc, minioClient.GetImgBucket(), model.Video)
	}
	return r.db.Model(&ds.Model{}).Where("model_id = ?", id).Update("is_deleted", true).Error
}
