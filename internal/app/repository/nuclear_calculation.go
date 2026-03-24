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

func (r *Repository) GetCalcModelCount(creatorID uint) int64 {
	var calcID uint
	err := r.db.Model(&ds.Nuclear_calculation{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("calc_id").First(&calcID).Error
	if err != nil {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.ModelCalc{}).
		Where("calc_id = ?", calcID).Count(&count).Error
	if err != nil {
		logrus.Error("error counting modelcalcs:", err)
	}
	return count
}

func (r *Repository) GetActiveCalcID(creatorID uint) uint {
	var calcID uint
	err := r.db.Model(&ds.Nuclear_calculation{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("calc_id").First(&calcID).Error
	if err != nil {
		return 0
	}
	return calcID
}

func (r *Repository) GetCalc(id int, creatorID uint) ([]ds.ModelCalc, *ds.Nuclear_calculation, error) {
	var calc ds.Nuclear_calculation
	err := r.db.Where("calc_id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&calc).Error
	if err != nil {
		return nil, nil, err
	}

	var items []ds.ModelCalc
	err = r.db.Where("calc_id = ?", id).
		Preload("Model").
		Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	return items, &calc, nil
}

func (r *Repository) AddModel(modelID uint, creatorID uint) error {
	var calc ds.Nuclear_calculation

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&calc).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		calc = ds.Nuclear_calculation{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&calc).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.ModelCalc{}).
		Where("calc_id = ? AND model_id = ?", calc.CalcID, modelID).
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

		item := ds.ModelCalc{
			CalcID:   calc.CalcID,
			ModelID:  modelID,
			Amount:   1,
			ResFuel:  &rf,
			ResPower: &rp,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteCalc(calcID uint) error {
	query := `
		UPDATE nuclear_calculations
		SET status = 'deleted'
		WHERE calc_id = $1;
	`
	result := r.db.Exec(query, calcID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("calculation with id %d not found", calcID)
	}
	return nil
}

func (r *Repository) IsDraftCalc(calcID int, creatorID uint) (bool, error) {
	var calc ds.Nuclear_calculation
	err := r.db.Select("status").Where("calc_id = ? AND creator_id = ?",
		calcID, creatorID).First(&calc).Error
	if err != nil {
		return false, err
	}
	return calc.Status == "draft", nil
}
