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
	ShortDesc   string
	ImageURL    string
	Video       string
	Power       float64
	FuelUsage   float64
}

type NuclearCalculation struct {
	ID          int
	Name        string
	Description string
	Models      []ModelCalc
	ModelCount  int
}

type ModelCalc struct {
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
			ShortDesc:   "Russian PWR powering icebreakers and future floating or small land-based nuclear plants",
			Power:       55,
			FuelUsage:   45,
			ImageURL:    "http://localhost:9000/reactorservice/ritm.jpg",
			Video:       "ritm.mp4",
		},
		{
			ID:          2,
			Name:        "HTR-PM",
			Description: "HTR-PM - Китайский малый модульный ядерный реактор. Это высокотемпературный газоохлаждаемый реактор четвертого поколения с шаровым топливом, разработанный на основе прототипа HTR-10.",
			ShortDesc:   "Chinese Gen-4 HTGR with pebble-bed fuel, developed as an advanced successor to the HTR-10",
			Power:       210,
			FuelUsage:   260,
			ImageURL:    "http://localhost:9000/reactorservice/htr.jpg",
			Video:       "htr.mp4",
		},
		{
			ID:          3,
			Name:        "КЛТ-40С",
			Description: "КЛТ-40С - Российская плавучая атомная теплоэлектростанция (ПАТЭС) проекта 20870, находящаяся в порту города Певек (Чаунский район, Чукотского автономного округа), самая северная АЭС в мире.",
			ShortDesc:   "Russian floating NPP deployed in Pevek, operating as the world’s northernmost nuclear facility",
			Power:       70,
			FuelUsage:   85,
			ImageURL:    "http://localhost:9000/reactorservice/klt.jpg",
			Video:       "klt.mp4",
		},
		{
			ID:          4,
			Name:        "IRIS",
			Description: "IRIS - Проект реактора четвертого поколения, разработанный международной командой компаний, лабораторий и университетов при координации компании Westinghouse, призван открыть новые рынки для атомной энергетики и создать мост между технологиями реакторов третьего и четвертого поколений.",
			ShortDesc:   "Westinghouse-led international Gen-4 PWR project bridging third and fourth-generation designs",
			Power:       335,
			FuelUsage:   320,
			ImageURL:    "http://localhost:9000/reactorservice/iris.jpg",
			Video:       "iris.mp4",
		},
		{
			ID:          5,
			Name:        "ГТ-МГР",
			Description: "ГТ-МГР - Российско-американский проект по созданию АЭС на базе высокотемпературного газоохлаждаемого реактора с гелиевым теплоносителем, работающего в прямом газотурбинном цикле.",
			ShortDesc:   "Russian-US helium-cooled high-temp reactor utilizing a direct gas-turbine power generation cycle",
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

func (r *Repository) buildCalc(id int, name, description string, entries []struct {
	ModelID int
	Amount  int
}) (NuclearCalculation, error) {
	models, err := r.GetModels()
	if err != nil {
		return NuclearCalculation{}, err
	}

	modMap := make(map[int]Model)
	for _, m := range models {
		modMap[m.ID] = m
	}

	var loadModels []ModelCalc

	for _, e := range entries {
		mod, ok := modMap[e.ModelID]
		if !ok {
			continue
		}
		res_p := CalculatePower(mod.Power, e.Amount)
		res_f := CalculateFuel(mod.FuelUsage, e.Amount)
		loadModels = append(loadModels, ModelCalc{
			Model:    mod,
			Amount:   e.Amount,
			ResPower: res_p,
			ResFuel:  res_f,
		})
	}

	return NuclearCalculation{
		ID:          id,
		Name:        name,
		Description: description,
		Models:      loadModels,
		ModelCount:  len(loadModels),
	}, nil
}

func (r *Repository) GetCalcs() ([]NuclearCalculation, error) {
	entries := []struct {
		ModelID int
		Amount  int
	}{
		{1, 3},
		{2, 2},
		{3, 1},
	}

	calc, err := r.buildCalc(
		1,
		"Тестовый расчёт",
		"Тестовый расчёт характеристик малых ядерных реакторов, показывает основной функционал калькулятора.",
		entries,
	)
	if err != nil {
		return nil, err
	}

	return []NuclearCalculation{calc}, nil
}

func (r *Repository) GetCalc(id int) (NuclearCalculation, error) {
	calcs, err := r.GetCalcs()
	if err != nil {
		return NuclearCalculation{}, err
	}

	for _, calc := range calcs {
		if calc.ID == id {
			return calc, nil
		}
	}
	return NuclearCalculation{}, fmt.Errorf("заявка не найдена")
}

func (r *Repository) GetCalcForModel(modelID int) (*ModelCalc, error) {
	calcs, err := r.GetCalcs()
	if err != nil {
		return nil, err
	}

	for _, calc := range calcs {
		for _, ms := range calc.Models {
			if ms.Model.ID == modelID {
				return &ms, nil
			}
		}
	}
	return nil, fmt.Errorf("модель не найдена в заявках")
}
