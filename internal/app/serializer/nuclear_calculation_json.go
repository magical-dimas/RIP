package serializer

import (
	"time"

	"rip_project/internal/app/ds"
)

type CalcJSON struct {
	CalcID             uint       `json:"calc_id"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatorLogin       string     `json:"creator_login"`
	ModeratorLogin     *string    `json:"moderator_login"`
	FormingDate        *time.Time `json:"forming_date"`
	FinishDate         *time.Time `json:"finish_date"`
	Description        *string    `json:"description"`
	CompletedItemCount int        `json:"completed_item_count"`
}

func CalcToJSON(calc ds.Nuclear_calculation, creatorLogin, moderatorLogin string, completedItemCount int) CalcJSON {
	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}
	var finishDate *time.Time
	if calc.FinishDate.Valid {
		finishDate = &calc.FinishDate.Time
	}
	return CalcJSON{
		CalcID:             calc.CalcID,
		Status:             calc.Status,
		CreatedAt:          calc.CreatedAt,
		CreatorLogin:       creatorLogin,
		ModeratorLogin:     mLogin,
		FormingDate:        calc.FormingDate,
		FinishDate:         finishDate,
		Description:        calc.Description,
		CompletedItemCount: completedItemCount,
	}
}

func CalcFromJSON(j CalcJSON) ds.Nuclear_calculation {
	return ds.Nuclear_calculation{
		Description: j.Description,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
