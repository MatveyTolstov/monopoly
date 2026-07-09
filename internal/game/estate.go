package game

type GroupColor int

const (
	Red GroupColor = iota
	Blue
	Brown
	Yellow
	Pink
	Orange
	Green
	DarkBlue
	Railroad
	Utility
)

type Estate struct {
	Name        string
	Cost        int    // цена покупки
	Rents       [5]int // рента домов 0, 1-3 дома, отель
	Group       GroupColor // цвет клеток для монополии
	HousePrice  int
	Owner       string // пока хз буду передовать имя игрока
	Houses      int
	IsMortgaged bool
}

func NewEstate(name string, cost int, rents [6]int, group GroupColor, housePrice int) *Estate {
	return &Estate{
		Name:       name,
		Cost:       cost,
		Rents:      rents,
		Group:      group,
		HousePrice: housePrice,
		Owner:      "",
		Houses:     0,
	}
}
