package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Model struct {
	ID          int
	Name        string
	Description string
	ImageURL    string
	Video       string
	Power       float64
	FuelUsage   float64
}

type Request struct {
	ID          int
	Name        string
	Description string
	Models      []ModelRequest
	ModelCount  int
}

type ModelRequest struct {
	Model    Model
	Amount   int
	ResPower float64
	ResFuel  float64
}

func (r *Repository) GetModels() ([]Model, error) {
	models := []Model{
		{
			ID:          1,
			Name:        "РИТМ-200",
			Description: "РИТМ-200 - водо-водяной ядерный реактор, предназначенный для установки на ледоколах и перспективных плавучих атомных электростанциях, малых АЭС.",
			Power:       55,
			FuelUsage:   45,
			ImageURL:    "http://localhost:9000/reactorservice/ritm.jpg",
			Video:       "ritm.mp4",
		},
		{
			ID:          2,
			Name:        "HTR-PM",
			Description: "HTR-PM - Китайский малый модульный ядерный реактор. Это высокотемпературный газоохлаждаемый реактор четвертого поколения с шаровым топливом, разработанный на основе прототипа HTR-10.",
			Power:       210,
			FuelUsage:   260,
			ImageURL:    "http://localhost:9000/reactorservice/htr.jpg",
			Video:       "htr.mp4",
		},
		{
			ID:          3,
			Name:        "КЛТ-40С",
			Description: "КЛТ-40С - Российская плавучая атомная теплоэлектростанция (ПАТЭС) проекта 20870, находящаяся в порту города Певек (Чаунский район, Чукотского автономного округа), самая северная АЭС в мире.",
			Power:       70,
			FuelUsage:   85,
			ImageURL:    "http://localhost:9000/reactorservice/klt.jpg",
			Video:       "klt.mp4",
		},
		{
			ID:          4,
			Name:        "IRIS",
			Description: "IRIS - Проект реактора четвертого поколения, разработанный международной командой компаний, лабораторий и университетов при координации компании Westinghouse, призван открыть новые рынки для атомной энергетики и создать мост между технологиями реакторов третьего и четвертого поколений.",
			Power:       335,
			FuelUsage:   320,
			ImageURL:    "http://localhost:9000/reactorservice/iris.jpg",
			Video:       "iris.mp4",
		},
		{
			ID:          5,
			Name:        "ГТ-МГР",
			Description: "ГТ-МГР - Российско-американский проект по созданию АЭС на базе высокотемпературного газоохлаждаемого реактора с гелиевым теплоносителем, работающего в прямом газотурбинном цикле.",
			Power:       285,
			FuelUsage:   300,
			ImageURL:    "http://localhost:9000/reactorservice/gtmgr.jpg",
			Video:       "gtmgr.mp4",
		},
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("массив моделей пуст")
	}

	return models, nil
}

func (r *Repository) GetModel(id int) (Model, error) {
	models, err := r.GetModels()
	if err != nil {
		return Model{}, err
	}

	for _, m := range models {
		if m.ID == id {
			return m, nil
		}
	}
	return Model{}, fmt.Errorf("модель не найдена")
}

func (r *Repository) GetModelByTitle(query string) ([]Model, error) {
	models, err := r.GetModels()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var result []Model
	for _, m := range models {
		if strings.Contains(strings.ToLower(m.Name), q) {
			result = append(result, m)
		}
	}
	return result, nil
}

func CalculatePower(power float64, amount int) float64 {
	return float64(amount) * power * 30
}

func CalculateFuel(fuel_usage float64, amount int) float64 {
	return float64(amount) * fuel_usage * 30
}

func (r *Repository) buildRequest(id int, name, description string, entries []struct {
	ModelID int
	Amount  int
}) (Request, error) {
	models, err := r.GetModels()
	if err != nil {
		return Request{}, err
	}

	modMap := make(map[int]Model)
	for _, m := range models {
		modMap[m.ID] = m
	}

	var loadModels []ModelRequest

	for _, e := range entries {
		mod, ok := modMap[e.ModelID]
		if !ok {
			continue
		}
		res_p := CalculatePower(mod.Power, e.Amount)
		res_f := CalculateFuel(mod.FuelUsage, e.Amount)
		loadModels = append(loadModels, ModelRequest{
			Model:    mod,
			Amount:   e.Amount,
			ResPower: res_p,
			ResFuel:  res_f,
		})
	}

	return Request{
		ID:          id,
		Name:        name,
		Description: description,
		Models:      loadModels,
		ModelCount:  len(loadModels),
	}, nil
}

func (r *Repository) GetRequests() ([]Request, error) {
	entries := []struct {
		ModelID int
		Amount  int
	}{
		{1, 3},
		{2, 2},
		{3, 1},
	}

	req, err := r.buildRequest(
		1,
		"Тестовый расчёт",
		"Тестовый расчёт характеристик малых ядерных реакторов, показывает основной функционал калькулятора.",
		entries,
	)
	if err != nil {
		return nil, err
	}

	return []Request{req}, nil
}

func (r *Repository) GetRequest(id int) (Request, error) {
	requests, err := r.GetRequests()
	if err != nil {
		return Request{}, err
	}

	for _, req := range requests {
		if req.ID == id {
			return req, nil
		}
	}
	return Request{}, fmt.Errorf("заявка не найдена")
}

func (r *Repository) GetRequestForModel(modelID int) (*ModelRequest, error) {
	requests, err := r.GetRequests()
	if err != nil {
		return nil, err
	}

	for _, req := range requests {
		for _, ms := range req.Models {
			if ms.Model.ID == modelID {
				return &ms, nil
			}
		}
	}
	return nil, fmt.Errorf("модель не найдена в заявках")
}
