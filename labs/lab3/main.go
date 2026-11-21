package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Damix222/tech-prog-2026/labs/lab3/rpg"
)

func main() {
	name := flag.String("name", "Герой", "имя нового героя")
	class := flag.String("class", "warrior", "warrior или ranger")
	seed := flag.Int64("seed", 0, "сид случайности; 0 — случайный")
	flag.Parse()
	if *seed != 0 {
		rpg.Dice().Seed(*seed)
	}
	var kit rpg.EquipmentFactory = rpg.WarriorKit{}
	if *class == "ranger" {
		kit = rpg.RangerKit{}
	} else if *class != "warrior" {
		fmt.Fprintln(os.Stderr, "неизвестный класс:", *class)
		os.Exit(1)
	}
	h, err := rpg.NewHeroBuilder().Named(*name).WithKit(kit).Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	locations := map[string]rpg.Location{"лес": rpg.Forest{}, "руины": rpg.Ruins{}, "дракон": rpg.DragonLair{}}
	fmt.Println("Команды: герой, лес, руины, дракон, отдых, вещи, надеть N, выход")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		parts := strings.Fields(strings.ToLower(scanner.Text()))
		if len(parts) == 0 {
			continue
		}
		err = nil
		switch parts[0] {
		case "герой":
			fmt.Printf("%s (%s): уровень %d, HP %d/%d, опыт %d/%d, золото %d\nОружие: %s (%d); броня: %s (%d)\n",
				h.Name, h.Class, h.Level, h.HP, h.MaxHP(), h.XP, h.Level*20, h.Gold, h.Weapon.Name, h.Weapon.Damage, h.Armor.Name, h.Armor.Defense)
		case "лес", "руины", "дракон":
			_, err = rpg.Explore(h, locations[parts[0]], os.Stdout)
		case "отдых":
			h.Rest()
			fmt.Println("Здоровье восстановлено.")
		case "вещи":
			for i, item := range h.Inventory {
				fmt.Printf("%d. %s: урон %d, руны %v\n", i+1, item.Name, item.Damage, item.Runes)
			}
		case "надеть":
			if len(parts) != 2 {
				err = fmt.Errorf("используйте: надеть N")
				break
			}
			index, parseErr := strconv.Atoi(parts[1])
			err = parseErr
			if err == nil {
				err = h.Equip(index - 1)
			}
		case "выход":
			return
		default:
			err = fmt.Errorf("неизвестная команда")
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
