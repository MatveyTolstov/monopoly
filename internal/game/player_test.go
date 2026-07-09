package game

import (
	"testing"
)

func newTestEstate(cost int) *Estate {
	return NewEstate("Тест", cost, [5]int{10, 20, 30, 40, 50}, Brown, 100)
}

// --- Pay ---

func TestPay_OK(t *testing.T) {
	p := NewPlayer(1, "Вася", 500)
	if err := p.Pay(200); err != nil {
		t.Fatal(err)
	}
	if p.Money != 300 {
		t.Fatalf("хотел 300, получил %d", p.Money)
	}
}

func TestPay_Exact(t *testing.T) {
	p := NewPlayer(1, "Вася", 200)
	if err := p.Pay(200); err != nil {
		t.Fatal(err)
	}
	if p.Money != 0 {
		t.Fatalf("хотел 0, получил %d", p.Money)
	}
}

func TestPay_InsufficientFunds(t *testing.T) {
	p := NewPlayer(1, "Вася", 100)
	if err := p.Pay(200); err == nil {
		t.Fatal("ожидалась ошибка, но её не было")
	}
	if p.Money != 100 {
		t.Fatalf("деньги не должны были измениться, получил %d", p.Money)
	}
}

// --- Move ---

func TestMove_Normal(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 5
	passed := p.Move(10)
	if p.Position != 15 {
		t.Fatalf("хотел 15, получил %d", p.Position)
	}
	if passed {
		t.Fatal("не должен был пройти СТАРТ")
	}
}

func TestMove_PassGo(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 35
	passed := p.Move(10)
	if p.Position != 5 {
		t.Fatalf("хотел 5, получил %d", p.Position)
	}
	if !passed {
		t.Fatal("должен был пройти СТАРТ")
	}
}

func TestMove_LandOnGo(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 35
	passed := p.Move(5)
	if p.Position != 0 {
		t.Fatalf("хотел 0 (СТАРТ), получил %d", p.Position)
	}
	if !passed {
		t.Fatal("приземление на СТАРТ должно считаться passedGo")
	}
}

func TestMove_NegativeSteps(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	p.Position = 2
	passed := p.Move(-5)
	if p.Position != 37 {
		t.Fatalf("хотел 37, получил %d", p.Position)
	}
	if passed {
		t.Fatal("движение назад не должно давать passedGo")
	}
}

// --- GoToJail / TickJail ---

func TestGoToJail(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	if err := p.GoToJail(); err != nil {
		t.Fatal(err)
	}
	if !p.InJail {
		t.Fatal("игрок должен быть в тюрьме")
	}
	if p.Position != 10 {
		t.Fatalf("позиция должна быть 10, получил %d", p.Position)
	}
	if p.JailTurnsLeft != 3 {
		t.Fatalf("JailTurnsLeft должен быть 3, получил %d", p.JailTurnsLeft)
	}
}

func TestGoToJail_AlreadyInJail(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	_ = p.GoToJail()
	if err := p.GoToJail(); err == nil {
		t.Fatal("ожидалась ошибка при повторной посадке в тюрьму")
	}
}

func TestTickJail_Release(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	_ = p.GoToJail()
	p.TickJail()
	p.TickJail()
	p.TickJail()
	if p.InJail {
		t.Fatal("игрок должен был выйти из тюрьмы после 3 ходов")
	}
	if p.JailTurnsLeft != 0 {
		t.Fatalf("JailTurnsLeft должен быть 0, получил %d", p.JailTurnsLeft)
	}
}

func TestTickJail_NotInJail(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	before := p.JailTurnsLeft
	p.TickJail()
	if p.JailTurnsLeft != before {
		t.Fatal("TickJail не должен менять состояние если игрок не в тюрьме")
	}
}

// --- Mortgage ---

func TestMortgage_OK(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)
	p.Actives = append(p.Actives, e)

	if err := p.Mortgage(e); err != nil {
		t.Fatal(err)
	}
	if !e.IsMortgaged {
		t.Fatal("недвижимость должна быть заложена")
	}
	if p.Money != 100 {
		t.Fatalf("хотел 100 (200/2), получил %d", p.Money)
	}
}

func TestMortgage_AlreadyMortgaged(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)
	e.IsMortgaged = true
	p.Actives = append(p.Actives, e)

	if err := p.Mortgage(e); err == nil {
		t.Fatal("ожидалась ошибка при повторном залоге")
	}
}

func TestMortgage_NotOwned(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)

	if err := p.Mortgage(e); err == nil {
		t.Fatal("ожидалась ошибка — недвижимость не принадлежит игроку")
	}
}

// --- SellEstate ---

func TestSellEstate_OK(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)
	e.Owner = "Вася"
	e.Houses = 2
	p.Actives = append(p.Actives, e)

	if err := p.SellEstate(e); err != nil {
		t.Fatal(err)
	}
	if p.Money != 100 {
		t.Fatalf("хотел 100 (200/2), получил %d", p.Money)
	}
	if e.Owner != "" {
		t.Fatal("владелец должен быть сброшен")
	}
	if e.Houses != 0 {
		t.Fatal("дома должны быть сброшены")
	}
	if len(p.Actives) != 0 {
		t.Fatal("недвижимость должна быть убрана из списка")
	}
}

func TestSellEstate_Mortgaged(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)
	e.IsMortgaged = true
	p.Actives = append(p.Actives, e)

	if err := p.SellEstate(e); err != nil {
		t.Fatal(err)
	}
	if p.Money != 0 {
		t.Fatalf("заложенная недвижимость не должна давать деньги при продаже, получил %d", p.Money)
	}
	if e.IsMortgaged {
		t.Fatal("IsMortgaged должен быть сброшен после продажи")
	}
}

func TestSellEstate_NotOwned(t *testing.T) {
	p := NewPlayer(1, "Вася", 0)
	e := newTestEstate(200)

	if err := p.SellEstate(e); err == nil {
		t.Fatal("ожидалась ошибка — недвижимость не принадлежит игроку")
	}
}
