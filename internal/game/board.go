package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type SpaceType int

const (
	SpaceGo           SpaceType = iota // старт, игрок получает деньги при прохождении
	SpaceProperty                      // обычная улица, можно купить и строить дома
	SpaceChance                        // карточка шанс, тянешь из колоды
	SpaceJustVisiting                  // просто посещение тюрьмы, ничего не происходит
	SpaceGoToJail                      // отправляет игрока в тюрьму
	SpaceFreeParking                   // бесплатная парковка, ничего не происходит
	SpaceRailroad                      // железнодорожная станция
	Shop24                             // магазин 24
)

var spaceTypeMap = map[string]SpaceType{
	"Go":           SpaceGo,
	"Property":     SpaceProperty,
	"Chance":       SpaceChance,
	"JustVisiting": SpaceJustVisiting,
	"GoToJail":     SpaceGoToJail,
	"FreeParking":  SpaceFreeParking,
	"Railroad":     SpaceRailroad,
	"Shop24":       Shop24,
}

var groupColorMap = map[string]GroupColor{
	"Red":      Red,
	"Blue":     Blue,
	"Brown":    Brown,
	"Yellow":   Yellow,
	"Pink":     Pink,
	"Orange":   Orange,
	"Green":    Green,
	"DarkBlue": DarkBlue,
	"Railroad": Railroad,
	"Utility":  Utility,
}

type estateJSON struct {
	Cost       int    `json:"cost"`
	Rents      [5]int `json:"rents"`
	Group      string `json:"group"`
	HousePrice int    `json:"house_price"`
}

type spaceJSON struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Estate *estateJSON `json:"estate,omitempty"`
}

type Space struct {
	Name   string
	Type   SpaceType
	Estate *Estate // nil для Go, Jail, Chance и т.д.
}

type Board struct {
	Spaces [40]Space
}

func NewBoard(spacesPath string) (*Board, error) {
	data, err := os.ReadFile(spacesPath)
	if err != nil {
		return nil, fmt.Errorf("не могу прочитать файл: %w", err)
	}

	var raw []spaceJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if len(raw) != 40 {
		return nil, fmt.Errorf("в JSON %d клеток, нужно 40", len(raw))
	}

	b := &Board{}
	for i, s := range raw {
		spaceType, ok := spaceTypeMap[s.Type]
		if !ok {
			return nil, fmt.Errorf("неизвестный тип клетки: %s", s.Type)
		}

		var estate *Estate
		if s.Estate != nil {
			group, ok := groupColorMap[s.Estate.Group]
			if !ok {
				return nil, fmt.Errorf("неизвестный цвет группы: %s", s.Estate.Group)
			}
			estate = NewEstate(s.Name, s.Estate.Cost, s.Estate.Rents, group, s.Estate.HousePrice)
		}

		b.Spaces[i] = Space{Name: s.Name, Type: spaceType, Estate: estate}
	}

	return b, nil
}
