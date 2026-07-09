package game

import (
	"errors"
	"slices"
)

type Player struct {
	ID            int
	Name          string
	Money         int
	Position      int
	JailTurnsLeft int
	InJail        bool
	Actives       []*Estate
	Bankrupt      bool
}

func NewPlayer(id int, name string, money int) *Player {
	return &Player{
		ID:      id,
		Name:    name,
		Money:   money,
		InJail:  false,
		Actives: make([]*Estate, 0),
	}
}

func (p *Player) Pay(amount int) error {
	if p.Money < amount {
		return errors.New("у игрока недостаточно денег")
	}

	p.Money -= amount

	return nil
}

func (p *Player) SellEstate(estate *Estate) error {
	i := slices.Index(p.Actives, estate)

	if i == -1 {
		return errors.New("этой недвижимости нет у игрока")
	}

	if !estate.IsMortgaged {
		p.Money += estate.Cost / 2
	}

	estate.Owner = ""
	estate.IsMortgaged = false
	estate.Houses = 0
	p.Actives = append(p.Actives[:i], p.Actives[i+1:]...)
	return nil
}

func (p *Player) Mortgage(estate *Estate) error {
	i := slices.Index(p.Actives, estate)

	if i == -1 {
		return errors.New("этой недвижимости нет у игрока")
	}

	if estate.IsMortgaged {
		return errors.New("недвижимость уже заложена")
	}

	p.Money += estate.Cost / 2

	estate.IsMortgaged = true

	return nil
}

func (p *Player) Move(steps int) bool {
	newPos := p.Position + steps
	passedGo := newPos >= 40
	p.Position = ((newPos % 40) + 40) % 40
	return passedGo
}

func (p *Player) GoToJail() error {
	if p.InJail {
		return errors.New("игрок уже в тюрьме")
	}

	p.Position = 10
	p.JailTurnsLeft = 3
	p.InJail = true

	return nil
}

func (p *Player) TickJail() {
	if !p.InJail {
		return
	}

	p.JailTurnsLeft--

	if p.JailTurnsLeft == 0 {
		p.InJail = false
	}
}
