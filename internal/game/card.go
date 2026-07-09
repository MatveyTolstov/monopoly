package game

import "fmt"

// CardType — тип события на карте.
type CardType string

const (
	CardPayTax      CardType = "pay_tax"
	CardCollect     CardType = "collect_money"
	CardGoToJail    CardType = "go_to_jail"
	CardMoveSteps   CardType = "move_steps"
	CardMoveTo      CardType = "move_to"
)

// CardParams — параметры конкретной карты.
// Поля не нужные для данного типа просто остаются нулевыми
// и не сериализуются в JSON благодаря omitempty.
type CardParams struct {
	Amount int `json:"amount,omitempty"`
	Steps  int `json:"steps,omitempty"`
	Pos    int `json:"pos,omitempty"`
}

// Card — то, что хранится в колоде и уходит клиенту как данные.
type Card struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Type   CardType   `json:"type"`
	Params CardParams `json:"params,omitempty"`
}

// ActionResult — то, что уходит клиенту после применения карты,
// чтобы фронт мог отрисовать/анимировать произошедшее.
type ActionResult struct {
	Type     CardType `json:"type"`
	Amount   int      `json:"amount,omitempty"`
	Steps    int      `json:"steps,omitempty"`
	Pos      int      `json:"pos,omitempty"`
	PassedGo bool     `json:"passedGo,omitempty"`
}

// Effect — сигнатура обработчика конкретного типа карты.
type Effect func(p *Player, params CardParams) (*ActionResult, error)

// effects — реестр обработчиков. Чтобы добавить новую карту,
// достаточно добавить сюда новый case-функцию, Card и Apply не трогаем.
var effects = map[CardType]Effect{
	CardPayTax: func(p *Player, params CardParams) (*ActionResult, error) {
		if err := p.Pay(params.Amount); err != nil {
			return nil, err
		}
		return &ActionResult{Type: CardPayTax, Amount: params.Amount}, nil
	},

	CardCollect: func(p *Player, params CardParams) (*ActionResult, error) {
		p.Money += params.Amount
		return &ActionResult{Type: CardCollect, Amount: params.Amount}, nil
	},

	CardGoToJail: func(p *Player, params CardParams) (*ActionResult, error) {
		if err := p.GoToJail(); err != nil {
			return nil, err
		}
		return &ActionResult{Type: CardGoToJail, Pos: p.Position}, nil
	},

	CardMoveSteps: func(p *Player, params CardParams) (*ActionResult, error) {
		passedGo := p.Move(params.Steps)
		return &ActionResult{
			Type:     CardMoveSteps,
			Steps:    params.Steps,
			Pos:      p.Position,
			PassedGo: passedGo,
		}, nil
	},

	CardMoveTo: func(p *Player, params CardParams) (*ActionResult, error) {
		passedGo := params.Pos < p.Position
		p.Position = params.Pos
		return &ActionResult{
			Type:     CardMoveTo,
			Pos:      p.Position,
			PassedGo: passedGo,
		}, nil
	},
}

// Apply применяет эффект карты к игроку и возвращает результат
// для отправки клиенту.
func (c *Card) Apply(p *Player) (*ActionResult, error) {
	effect, ok := effects[c.Type]
	if !ok {
		return nil, fmt.Errorf("неизвестный тип карты: %s", c.Type)
	}
	return effect(p, c.Params)
}