package serializer

import "rip_project/internal/app/ds"

type ModelJSON struct {
	ModelID     uint    `json:"model_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ShortDesc   string  `json:"short_desc"`
	IsDeleted   bool    `json:"is_deleted"`
	PhotoURL    string  `json:"photo_url"`
	Video       string  `json:"video"`
	Power       float64 `json:"power"`
	FuelUsage   float64 `json:"fuel_usage"`
}

func ModelToJSON(m ds.Model) ModelJSON {
	return ModelJSON{
		ModelID:     m.ModelID,
		Title:       m.Title,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		IsDeleted:   m.IsDeleted,
		PhotoURL:    m.PhotoURL,
		Video:       m.Video,
		Power:       m.Power,
		FuelUsage:   m.FuelUsage,
	}
}

func ModelFromJSON(j ModelJSON) ds.Model {
	return ds.Model{
		Title:       j.Title,
		Description: j.Description,
		Power:       j.Power,
		FuelUsage:   j.FuelUsage,
	}
}
