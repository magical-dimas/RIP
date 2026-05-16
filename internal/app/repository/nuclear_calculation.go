package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rip_project/internal/app/ds"
	"rip_project/internal/app/serializer"
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

func (r *Repository) CheckCurrentDraft(creatorID uint) (ds.Nuclear_calculation, error) {
	if creatorID == 0 {
		return ds.Nuclear_calculation{}, ErrNotAllowed
	}
	var calc ds.Nuclear_calculation
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&calc)
	if res.Error != nil {
		return ds.Nuclear_calculation{}, res.Error
	}
	if res.RowsAffected == 0 {
		return ds.Nuclear_calculation{}, ErrNoDraft
	}
	return calc, nil
}

func (r *Repository) GetCalcDraft(creatorID uint) (ds.Nuclear_calculation, bool, error) {
	calc, err := r.CheckCurrentDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		calc = ds.Nuclear_calculation{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&calc).Error; err != nil {
			return ds.Nuclear_calculation{}, false, err
		}
		return calc, true, nil
	}
	if err != nil {
		return ds.Nuclear_calculation{}, false, err
	}
	return calc, false, nil
}

func (r *Repository) GetModeratorAndCreatorLogin(calc ds.Nuclear_calculation) (string, string, error) {
	var creator ds.Engineers
	if err := r.db.Where("engineer_id = ?", calc.CreatorID).First(&creator).Error; err != nil {
		return "", "", err
	}
	var moderatorLogin string
	if calc.ModeratorID != nil && *calc.ModeratorID != 0 {
		var moderator ds.Engineers
		if err := r.db.Where("engineer_id = ?", *calc.ModeratorID).First(&moderator).Error; err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) GetCompletedItemCount(calcID uint) (int, error) {
	var count int64
	var calc ds.Nuclear_calculation
	res := r.db.Where("calc_id = ?", calcID).Limit(1).Find(&calc)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, ErrNotFound
	}
	if (calc.Status != "completed" && calc.Status != "formed"){
		return 0, nil
	}
	err := r.db.Model(&ds.ModelCalc{}).
		Where("calc_id = ? AND res_power IS NOT NULL AND res_fuel IS NOT NULL", calcID).
		Count(&count).Error
	return int(count), err
}

func (r *Repository) GetAllCalcs(from, to time.Time, status string, uid uint, is_mod bool) ([]ds.Nuclear_calculation, error) {
	var calcs []ds.Nuclear_calculation
    sub := r.db.Where("status != ? AND status != ?", "deleted", "draft")
	if !is_mod {
        sub = sub.Where("creator_id = ?", uid)
    }
	if !from.IsZero() {
		sub = sub.Where("forming_date >= ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("forming_date < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("calc_id").Find(&calcs).Error
	return calcs, err
}

func (r *Repository) GetSingleCalc(id int) (ds.Nuclear_calculation, error) {
	if id < 0 {
		return ds.Nuclear_calculation{}, errors.New("неверное id")
	}
	var calc ds.Nuclear_calculation
	err := r.db.Where("calc_id = ?", id).First(&calc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Nuclear_calculation{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status == "deleted" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return calc, nil
}

func (r *Repository) GetCalcItems(calcID int) ([]ds.ModelCalc, error) {
	var items []ds.ModelCalc
	err := r.db.Where("calc_id = ?", calcID).
		Preload("Model").
		Find(&items).Error
	return items, err
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
	var model ds.Model
	if err := r.db.Where("model_id = ? AND is_deleted = ?", modelID, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: стратегия с id %d", ErrNotFound, modelID)
		}
		return err
	}

	var count int64
	r.db.Model(&ds.ModelCalc{}).
		Where("calc_id = ? AND model_id = ?", calc.CalcID, modelID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("%w: модель %d уже в заявке %d", ErrAlreadyExists, modelID, calc.CalcID)
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
	return r.db.Create(&item).Error
}

func (r *Repository) DeleteModelFromCalc(calcID, modelID int) (ds.Nuclear_calculation, error) {
	var calc ds.Nuclear_calculation
	err := r.db.Where("calc_id = ?", calcID).First(&calc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Nuclear_calculation{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, calcID)
		}
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status != "draft" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: можно удалять только из черновика", ErrNotAllowed)
	}
	err = r.db.Where("calc_id = ? AND model_id = ?", calcID, modelID).
		Delete(&ds.ModelCalc{}).Error
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	return calc, nil
}

func (r *Repository) EditModelInCalc(calcID, modelID int, j serializer.ModelCalcJSON) (ds.ModelCalc, error) {
	var item ds.ModelCalc
	err := r.db.Where("calc_id = ? AND model_id = ?", calcID, modelID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ModelCalc{}, fmt.Errorf("%w: модель в заявке", ErrNotFound)
		}
		return ds.ModelCalc{}, err
	}

	var calc ds.Nuclear_calculation
	if err := r.db.Where("calc_id = ?", calcID).First(&calc).Error; err != nil {
		return ds.ModelCalc{}, err
	}
	if calc.Status != "draft" {
		return ds.ModelCalc{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}

	updates := map[string]interface{}{
		"amount": j.Amount,
	}

	var model ds.Model
	if err := r.db.First(&model, modelID).Error; err == nil {
		rf := CalculateFuel(
			model.FuelUsage,
			j.Amount,
		)
		rp := CalculatePower(
			model.Power,
			j.Amount,
		)
		updates["res_fuel"] = rf
		updates["res_power"] = rp
	}

	err = r.db.Model(&item).Updates(updates).Error
	if err != nil {
		return ds.ModelCalc{}, err
	}
	r.db.Where("calc_id = ? AND model_id = ?", calcID, modelID).
		Preload("Model").First(&item)
	return item, nil
}

func (r *Repository) EditCalc(id int, j serializer.CalcJSON) (ds.Nuclear_calculation, error) {
	var calc ds.Nuclear_calculation
	err := r.db.Where("calc_id = ? AND status != ?", id, "deleted").First(&calc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Nuclear_calculation{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status != "draft" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}
	updates := serializer.CalcFromJSON(j)
	err = r.db.Model(&calc).Updates(updates).Error
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	r.db.Where("calc_id = ?", id).First(&calc)
	return calc, nil
}
func (r *Repository) FormCalc(id int) (ds.Nuclear_calculation, error) {
	calc, err := r.GetSingleCalc(id)
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status != "draft" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: только черновик можно сформировать", ErrNotAllowed)
	}

	items, err := r.GetCalcItems(int(calc.CalcID))
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	if len(items) == 0 {
		return ds.Nuclear_calculation{}, errors.New("нельзя сформировать пустую заявку")
	}

	for _, item := range items {
		var model ds.Model
		if err := r.db.First(&model, item.ModelID).Error; err != nil {
			return ds.Nuclear_calculation{}, err
		}
		rf := CalculateFuel(
			model.FuelUsage,
			item.Amount,
		)
		rp := CalculatePower(
			model.Power,
			item.Amount,
		)
		r.db.Model(&ds.ModelCalc{}).
			Where("calc_id = ? AND model_id = ?", calc.CalcID, item.ModelID).
			Update("res_fuel", rf)
		r.db.Model(&ds.ModelCalc{}).
			Where("calc_id = ? AND model_id = ?", calc.CalcID, item.ModelID).
			Update("res_power", rp)
	}

	formingDate := time.Now()
	err = r.db.Model(&calc).Updates(map[string]interface{}{
		"status":       "formed",
		"forming_date": formingDate,
	}).Error
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	calc.Status = "formed"
	calc.FormingDate = &formingDate
	return calc, nil
}
func (r *Repository) FinishCalc(id int, status string, userId int) (ds.Nuclear_calculation, error) {
	if status != "completed" && status != "rejected" {
		return ds.Nuclear_calculation{}, errors.New("неверный статус: допустимы completed или rejected")
	}
	engineer, err := r.GetEngineerByID(userId)
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	if !engineer.IsTechnician {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}
	calc, err := r.GetSingleCalc(id)
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status != "formed" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: завершить/отклонить можно только сформированную заявку", ErrNotAllowed)
	}
	finishDate := time.Now()
	err = r.db.Model(&calc).Updates(map[string]interface{}{
		"status":       status,
		"finish_date":  finishDate,
		"moderator_id": engineer.EngineerID,
	}).Error
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	calc.Status = status
	calc.FinishDate = &finishDate
	uid := engineer.EngineerID
	calc.ModeratorID = &uid
	return calc, nil
}

func (r *Repository) DeleteCalc(calcID int) (ds.Nuclear_calculation, error) {
	calc, err := r.GetSingleCalc(calcID)
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	if calc.Status != "draft" {
		return ds.Nuclear_calculation{}, fmt.Errorf("%w: удалить можно только черновик", ErrNotAllowed)
	}
	formingDate := time.Now()
	err = r.db.Model(&calc).Updates(map[string]interface{}{
		"status":       "deleted",
		"forming_date": formingDate,
	}).Error
	if err != nil {
		return ds.Nuclear_calculation{}, err
	}
	calc.Status = "deleted"
	calc.FormingDate = &formingDate
	return calc, nil
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
