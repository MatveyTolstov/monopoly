package game

import (
	"testing"
)

func makeCard(t CardType, params CardParams) *Card {
	return &Card{ID: "test", Name: "тест", Type: t, Params: params}
}

func TestCard_PayTax(t *testing.T) {
	p := NewPlayer(1, "Вася", 500)
	_, err := makeCard(CardPayTax, CardParams{Amount: 100}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Money != 400 {
		t.Fatalf("хотел 400, получил %d", p.Money)
	}
}

func TestCard_PayTax_Broke(t *testing.T) {
	p := NewPlayer(1, "Вася", 50)
	_, err := makeCard(CardPayTax, CardParams{Amount: 100}).Apply(p)
	if err == nil {
		t.Fatal("ожидалась ошибка при нехватке денег")
	}
}

func TestCard_Collect(t *testing.T) {
	p := NewPlayer(1, "Вася", 100)
	_, err := makeCard(CardCollect, CardParams{Amount: 200}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Money != 300 {
		t.Fatalf("хотел 300, получил %d", p.Money)
	}
}

func TestCard_GoToJail(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 25
	_, err := makeCard(CardGoToJail, CardParams{}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if !p.InJail {
		t.Fatal("игрок должен быть в тюрьме")
	}
	if p.Position != 10 {
		t.Fatalf("хотел позицию 10, получил %d", p.Position)
	}
}

func TestCard_MoveSteps(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 10
	res, err := makeCard(CardMoveSteps, CardParams{Steps: 5}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Position != 15 {
		t.Fatalf("хотел 15, получил %d", p.Position)
	}
	if res.PassedGo {
		t.Fatal("не должен был пройти СТАРТ")
	}
}

func TestCard_MoveSteps_PassGo(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 38
	res, err := makeCard(CardMoveSteps, CardParams{Steps: 5}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Position != 3 {
		t.Fatalf("хотел 3, получил %d", p.Position)
	}
	if !res.PassedGo {
		t.Fatal("должен был пройти СТАРТ")
	}
}

func TestCard_MoveTo(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 20
	res, err := makeCard(CardMoveTo, CardParams{Pos: 5}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Position != 5 {
		t.Fatalf("хотел 5, получил %d", p.Position)
	}
	if !res.PassedGo {
		t.Fatal("должен был пройти СТАРТ (5 < 20)")
	}
}

func TestCard_MoveTo_NoPassGo(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 5
	res, err := makeCard(CardMoveTo, CardParams{Pos: 20}).Apply(p)
	if err != nil {
		t.Fatal(err)
	}
	if res.PassedGo {
		t.Fatal("не должен был пройти СТАРТ (20 > 5)")
	}
}

func TestCard_UnknownType(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	c := &Card{Type: "несуществующий_тип"}
	_, err := c.Apply(p)
	if err == nil {
		t.Fatal("ожидалась ошибка для неизвестного типа карты")
	}
}
