package rpg

import (
	"fmt"
	"io"
)

func (h *Hero) Rest() { h.HP = h.MaxHP() }

func (h *Hero) Equip(index int) error {
	if index < 0 || index >= len(h.Inventory) {
		return fmt.Errorf("нет предмета с таким номером")
	}
	h.Weapon, h.Inventory[index] = h.Inventory[index], h.Weapon
	return nil
}

// Explore работает с интерфейсами Location и Enemy, а тип врага выбирает локация.
func Explore(h *Hero, location Location, out io.Writer) (bool, error) {
	if h.HP <= 0 {
		return false, fmt.Errorf("сначала отдохните в лагере")
	}
	enemy := location.CreateEnemy()
	if _, boss := enemy.(*Dragon); boss && h.Level < 3 {
		return false, fmt.Errorf("для дракона нужен хотя бы 3-й уровень")
	}
	fmt.Fprintf(out, "%s: %s против %s\n", location.Name(), h.Name, enemy.Name())
	for round := 1; h.HP > 0 && enemy.Health() > 0 && round <= 100; round++ {
		dealt := enemy.Receive(h.Weapon.RollDamage() + h.Level - 1)
		taken := 0
		if enemy.Health() > 0 {
			taken = min(h.HP, max(0, enemy.Attack()-h.Armor.Defense))
			h.HP -= taken
		}
		fmt.Fprintf(out, "  Ход %d: урон %d/%d; HP героя %d, врага %d\n", round, dealt, taken, h.HP, enemy.Health())
	}
	won := enemy.Health() == 0
	if !won {
		h.Gold /= 2
		fmt.Fprintln(out, "Экспедиция провалена. Потеряна половина золота; можно отдохнуть в лагере.")
		return false, nil
	}
	h.Gold += location.Reward()
	h.XP += location.Reward()
	for h.XP >= h.Level*20 {
		h.XP -= h.Level * 20
		h.Level++
		h.Rest()
	}
	loot := h.Weapon.Clone()
	loot.Name, loot.Damage = "Трофей из локации «"+location.Name()+"»", h.Weapon.Damage+1
	loot.Runes = append(loot.Runes, "трофейная руна")
	h.Inventory = append(h.Inventory, loot)
	fmt.Fprintf(out, "Победа! +%d золота и опыта; добыча: %s\n", location.Reward(), loot.Name)
	if _, boss := enemy.(*Dragon); boss {
		fmt.Fprintln(out, "Дракон побеждён — цель игры выполнена!")
	}
	return true, nil
}
