package serializer

import "rip_project/internal/app/ds"

type ModelCalcJSON struct {
	CalcID   uint     `json:"calc_id"`
	ModelID  uint     `json:"model_id"`
	Amount   int      `json:"amount"`
	ResFuel  *float64 `json:"res_fuel"`
	ResPower *float64 `json:"res_power"`
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
	}
}

func ModelCalcDetailToJSON(item ds.ModelCalc) ModelCalcDetailJSON {
	return ModelCalcDetailJSON{
		CalcID:   item.CalcID,
		ModelID:  item.ModelID,
		Amount:   item.Amount,
		ResFuel:  item.ResFuel,
		ResPower: item.ResPower,
		Model:    ModelToJSON(item.Model),
	}
}

func ModelCalcFromJSON(j ModelCalcJSON) ds.ModelCalc {
	return ds.ModelCalc{
		Amount: j.Amount,
	}
}
