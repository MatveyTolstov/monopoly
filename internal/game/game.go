package game

type Game struct {
	Board              *Board
	Players            []*Player
	CurrentPlayerIndex int
}
