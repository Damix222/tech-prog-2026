package rpg

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"
)

// Singleton: тип закрыт для других пакетов, создание проходит через sync.Once.
type dice struct {
	mu  sync.Mutex
	rng *rand.Rand
}

var diceOnce sync.Once
var sharedDice *dice

type Random interface {
	Seed(int64)
	Roll(int) int
}

func Dice() Random {
	diceOnce.Do(func() { sharedDice = &dice{rng: rand.New(rand.NewSource(time.Now().UnixNano()))} })
	return sharedDice
}

func (d *dice) Seed(seed int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rng.Seed(seed)
}

func (d *dice) Roll(sides int) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.rng.Intn(sides) + 1
}

type Weapon struct {
	Name   string
	Damage int
	Runes  []string
}

// Prototype: изменяемый слайс копируется отдельно от структуры.
func (w Weapon) Clone() Weapon {
	w.Runes = slices.Clone(w.Runes)
	return w
}

func (w Weapon) RollDamage() int { return w.Damage + Dice().Roll(4) - 1 }

type Armor struct {
	Name    string
	Defense int
}

// Abstract Factory: каждая реализация создает согласованную пару предметов.
type EquipmentFactory interface {
	CreateWeapon() Weapon
	CreateArmor() Armor
	Class() string
}

var sword = Weapon{Name: "Меч", Damage: 8, Runes: []string{"стойкость"}}
var bow = Weapon{Name: "Лук", Damage: 10, Runes: []string{"точность"}}

type WarriorKit struct{}

func (WarriorKit) Class() string        { return "Воин" }
func (WarriorKit) CreateWeapon() Weapon { return sword.Clone() }
func (WarriorKit) CreateArmor() Armor   { return Armor{"Кольчуга", 3} }

type RangerKit struct{}

func (RangerKit) Class() string        { return "Следопыт" }
func (RangerKit) CreateWeapon() Weapon { return bow.Clone() }
func (RangerKit) CreateArmor() Armor   { return Armor{"Кожаная броня", 1} }

type Hero struct {
	Name, Class         string
	HP, Level, XP, Gold int
	Weapon              Weapon
	Armor               Armor
	Inventory           []Weapon
}

func (h *Hero) MaxHP() int { return 40 + (h.Level-1)*10 }

func (h *Hero) clone() *Hero {
	c := *h
	c.Weapon = h.Weapon.Clone()
	c.Inventory = slices.Clone(h.Inventory)
	for i := range c.Inventory {
		c.Inventory[i] = c.Inventory[i].Clone()
	}
	return &c
}

func (h *Hero) validate() error {
	if strings.TrimSpace(h.Name) == "" || h.Class == "" || h.Level < 1 || h.Level > 1000 ||
		h.Gold < 0 || h.XP < 0 || h.XP >= h.Level*20 || h.HP < 0 || h.HP > h.MaxHP() ||
		h.Weapon.Damage <= 0 || h.Armor.Defense < 0 {
		return fmt.Errorf("некорректные характеристики героя")
	}
	return nil
}

// Builder: промежуточное состояние хранится отдельно, Build проверяет результат.
type HeroBuilder struct{ hero Hero }

func NewHeroBuilder() *HeroBuilder {
	return &HeroBuilder{hero: Hero{Level: 1, Gold: 10}}
}

func (b *HeroBuilder) Named(name string) *HeroBuilder {
	b.hero.Name = strings.TrimSpace(name)
	return b
}

func (b *HeroBuilder) WithKit(kit EquipmentFactory) *HeroBuilder {
	b.hero.Class = kit.Class()
	b.hero.Weapon, b.hero.Armor = kit.CreateWeapon(), kit.CreateArmor()
	return b
}

func (b *HeroBuilder) Build() (*Hero, error) {
	h := b.hero.clone()
	h.HP = h.MaxHP()
	if err := h.validate(); err != nil {
		return nil, err
	}
	return h, nil
}

type Enemy interface {
	Name() string
	Health() int
	Attack() int
	Receive(int) int
}

type monster struct {
	name             string
	hp, power, armor int
}

func (m *monster) Name() string { return m.name }
func (m *monster) Health() int  { return m.hp }
func (m *monster) Attack() int  { return m.power + Dice().Roll(3) - 1 }
func (m *monster) Receive(hit int) int {
	damage := max(0, hit-m.armor)
	damage = min(m.hp, damage)
	m.hp -= damage
	return damage
}

type Goblin struct{ monster }
type Skeleton struct{ monster }
type Dragon struct{ monster }

// Factory Method: общая экспедиция вызывает CreateEnemy у выбранной локации.
type Location interface {
	Name() string
	CreateEnemy() Enemy
	Reward() int
}

type Forest struct{}
type Ruins struct{}
type DragonLair struct{}

func (Forest) Name() string       { return "Лес" }
func (Forest) Reward() int        { return 20 }
func (Forest) CreateEnemy() Enemy { return &Goblin{monster{"Гоблин", 22, 5, 0}} }

func (Ruins) Name() string       { return "Руины" }
func (Ruins) Reward() int        { return 30 }
func (Ruins) CreateEnemy() Enemy { return &Skeleton{monster{"Скелет", 35, 7, 1}} }

func (DragonLair) Name() string       { return "Логово дракона" }
func (DragonLair) Reward() int        { return 80 }
func (DragonLair) CreateEnemy() Enemy { return &Dragon{monster{"Дракон", 80, 12, 3}} }
