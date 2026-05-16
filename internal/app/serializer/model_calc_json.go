package serializer

import "rip_project/internal/app/ds"

type ModelCalcJSON struct {
	CalcID   uint     `json:"calc_id"`
	ModelID  uint     `json:"model_id"`
	Amount   int      `json:"amount"`
	ResFuel  *float64 `json:"res_fuel"`
	ResPower *float64 `json:"res_power"`
    Fuel     float64  `json:"fuel"`
    Power    float64  `json:"power"`
    Photo    string   `json:"photo"`
    Title string `json:"title"`
}

type ModelCalcDetailJSON struct {
	CalcID   uint      `json:"calc_id"`
	ModelID  uint      `json:"model_id"`
	Amount   int       `json:"amount"`
	ResFuel  *float64  `json:"res_fuel"`
	ResPower *float64  `json:"res_power"`
	Model    ModelJSON `json:"model"`
}

func ModelCalcToJSON(item ds.ModelCalc) ModelCalcJSON {
	return ModelCalcJSON{
		CalcID:   item.CalcID,
		ModelID:  item.ModelID,
		Amount:   item.Amount,
		ResFuel:  item.ResFuel,
		ResPower: item.ResPower,
        Fuel:     item.Model.FuelUsage,
        Power:    item.Model.Power,
        Photo:    item.Model.PhotoURL,
        Title:    item.Model.Title,
	}
}

func ModelCalcFromJSON(j ModelCalcJSON) ds.ModelCalc {
	return ds.ModelCalc{
		Amount: j.Amount,
	}
}
