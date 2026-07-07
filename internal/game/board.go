package game

type SpaceType int

const (
	SpaceGo             SpaceType = iota // старт, игрок получает 200 при прохождении
	SpaceProperty                        // обычная улица, можно купить и строить дома
	SpaceChance                          // карточка шанс, тянешь из колоды
	SpaceCommunityChest                  // общественная казна, тянешь из колоды
	SpaceJustVisiting                    // просто посещение тюрьмы, ничего не происходит
	SpaceGoToJail                        // отправляет игрока в тюрьму
	SpaceFreeParking                     // бесплатная парковка, ничего не происходит
)

type Space struct {
	Index  int
	Name   string
	Type   SpaceType
	Estate *Estate // nil для Go, Jail, Chance и т.д.
}

type Board struct {
	Spaces [40]Space
}

