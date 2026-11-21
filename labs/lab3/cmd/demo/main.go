package main

import (
	"fmt"
	"github.com/Damix222/tech-prog-2026/labs/lab3/rpg"
	"os"
)

func main() {
	fmt.Println("ЛАБА 3 — ПОРОЖДАЮЩИЕ ПАТТЕРНЫ")
	a, b := rpg.Dice(), rpg.Dice()
	fmt.Printf("Singleton: %p == %p: %t\n", a, b, a == b)
	h, err := rpg.NewHeroBuilder().Named("Ярослав").WithKit(rpg.WarriorKit{}).Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Builder: %s, %s, HP=%d, оружие=%s\n", h.Name, h.Class, h.HP, h.Weapon.Name)
	for _, kit := range []rpg.EquipmentFactory{rpg.WarriorKit{}, rpg.RangerKit{}} {
		fmt.Printf("Abstract Factory: %s -> %s + %s\n", kit.Class(), kit.CreateWeapon().Name, kit.CreateArmor().Name)
	}
	for _, location := range []rpg.Location{rpg.Forest{}, rpg.Ruins{}} {
		enemy := location.CreateEnemy()
		fmt.Printf("Factory Method: %s -> %T, %s, HP=%d\n", location.Name(), enemy, enemy.Name(), enemy.Health())
	}
	original := h.Weapon.Clone()
	shallow := original
	shallow.Runes[0] = "огонь"
	fmt.Printf("Поверхностная копия: оригинал=%v, копия=%v (общий массив)\n", original.Runes, shallow.Runes)
	original = h.Weapon.Clone()
	deep := original.Clone()
	deep.Name, deep.Runes[0] = "Копия меча", "лёд"
	fmt.Printf("Prototype, глубокая копия: оригинал=%s %v; копия=%s %v\n", original.Name, original.Runes, deep.Name, deep.Runes)
}
