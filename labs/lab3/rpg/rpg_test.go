package rpg

import (
	"io"
	"reflect"
	"testing"
)

func testHero(t *testing.T) *Hero {
	t.Helper()
	h, err := NewHeroBuilder().Named("Тест").WithKit(WarriorKit{}).Build()
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestCreationAndIsolation(t *testing.T) {
	if Dice() != Dice() {
		t.Fatal("singleton вернул разные экземпляры")
	}
	if _, err := NewHeroBuilder().WithKit(WarriorKit{}).Build(); err == nil {
		t.Fatal("строитель принял пустое имя")
	}
	builder := NewHeroBuilder().Named("Первый").WithKit(WarriorKit{})
	first, _ := builder.Build()
	first.Weapon.Runes[0] = "изменено"
	second, _ := builder.Named("Второй").Build()
	if second.Name != "Второй" || first.Name != "Первый" || second.Weapon.Runes[0] == "изменено" {
		t.Fatal("результаты Build разделяют изменяемое состояние")
	}
	w := WarriorKit{}.CreateWeapon()
	clone := w.Clone()
	clone.Runes[0] = "огонь"
	if w.Runes[0] == clone.Runes[0] || (WarriorKit{}).CreateWeapon().Runes[0] == "огонь" {
		t.Fatal("клонирование повредило оригинал или шаблон фабрики")
	}
	if reflect.TypeOf(Forest{}.CreateEnemy()) == reflect.TypeOf(Ruins{}.CreateEnemy()) {
		t.Fatal("локации должны создавать разные типы врагов")
	}
}

func TestDamageAndCampaign(t *testing.T) {
	enemy := Ruins{}.CreateEnemy()
	if enemy.Receive(-5) != 0 || enemy.Health() != 35 || enemy.Receive(100) != 35 || enemy.Health() != 0 {
		t.Fatal("неверная обработка урона и здоровья")
	}
	Dice().Seed(23)
	h := testHero(t)
	if _, err := Explore(h, &DragonLair{}, io.Discard); err == nil {
		t.Fatal("дракон доступен до 3-го уровня")
	}
	if err := h.Equip(-1); err == nil {
		t.Fatal("принят неверный номер предмета")
	}
	for i := 0; i < 12; i++ {
		h.Rest()
		won, err := Explore(h, Forest{}, io.Discard)
		if err != nil || !won {
			t.Fatalf("не пройдена экспедиция: %v", err)
		}
		if err := h.Equip(len(h.Inventory) - 1); err != nil {
			t.Fatal(err)
		}
	}
	h.Rest()
	if won, err := Explore(h, DragonLair{}, io.Discard); err != nil || !won {
		t.Fatalf("дракон не побеждён: %v", err)
	}
	h.HP = 0
	if _, err := Explore(h, Forest{}, io.Discard); err == nil {
		t.Fatal("герой без здоровья отправился в экспедицию")
	}
}
