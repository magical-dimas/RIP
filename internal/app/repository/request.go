package repository

import (
	"errors"
	"fmt"
	"rip_project/internal/app/ds"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CalculatePower(power float64, amount int) float64 {
	return float64(amount) * power * 30
}

func CalculateFuel(fuel_usage float64, amount int) float64 {
	return float64(amount) * fuel_usage * 30
}

func (r *Repository) GetRequestModelCount(creatorID uint) int64 {
	var reqID uint
	err := r.db.Model(&ds.Request{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("request_id").First(&reqID).Error
	if err != nil {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.ModelRequest{}).
		Where("request_id = ?", reqID).Count(&count).Error
	if err != nil {
		logrus.Error("error counting modelrequests:", err)
	}
	return count
}

func (r *Repository) GetActiveRequestID(creatorID uint) uint {
	var reqID uint
	err := r.db.Model(&ds.Request{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("request_id").First(&reqID).Error
	if err != nil {
		return 0
	}
	return reqID
}

func (r *Repository) GetRequest(id int, creatorID uint) ([]ds.ModelRequest, *ds.Request, error) {
	var req ds.Request
	err := r.db.Where("request_id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&req).Error
	if err != nil {
		return nil, nil, err
	}

	var items []ds.ModelRequest
	err = r.db.Where("request_id = ?", id).
		Preload("Model").
		Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	return items, &req, nil
}

func (r *Repository) AddModel(modelID uint, creatorID uint) error {
	var req ds.Request

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&req).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		req = ds.Request{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&req).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.ModelRequest{}).
		Where("request_id = ? AND model_id = ?", req.RequestID, modelID).
		Count(&count)

	if count == 0 {
		var model ds.Model
		if err := r.db.First(&model, modelID).Error; err != nil {
			return err
		}

		rf := CalculateFuel(
			model.FuelUsage,
			1,
		)
		rp := CalculatePower(
			model.Power,
			1,
		)

		item := ds.ModelRequest{
			RequestID: req.RequestID,
			ModelID:   modelID,
			Amount:    1,
			ResFuel:   &rf,
			ResPower:  &rp,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteRequest(reqID uint) error {
	query := `
		UPDATE requests
		SET status = 'deleted'
		WHERE request_id = $1;
	`
	result := r.db.Exec(query, reqID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("request with id %d not found", reqID)
	}
	return nil
}

func (r *Repository) IsDraftRequest(reqID int, creatorID uint) (bool, error) {
	var req ds.Request
	err := r.db.Select("status").Where("request_id = ? AND creator_id = ?",
		reqID, creatorID).First(&req).Error
	if err != nil {
		return false, err
	}
	return req.Status == "draft", nil
}
