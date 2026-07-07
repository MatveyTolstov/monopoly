package game

import "math/rand"

type Player struct {
	ID        int
	Name      string
	Money     int
	Position  int
	JailTurns int
	InJail    bool
	Actives   []*Estate
	Bankrupt  bool
}

func NewPlayer(name string, money int) *Player {
	return &Player{
		ID:      rand.Int(),
		Name:    name,
		Money:   money,
		InJail:  false,
		Actives: make([]*Estate, 0),
	}
}
