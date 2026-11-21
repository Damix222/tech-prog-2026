package main

import "fmt"

type Pet interface {
	Express()
	GetName() string
	GetAge() int
	GetAgeInHumanYears() int
}

type Person struct {
	name string
	pet  Pet
}

type Animal struct {
	name string
	age  int
}

type Parrot struct {
	Animal
}

type Turtle struct {
	Animal
}

func (animal Animal) GetName() string {
	return animal.name
}

func (animal Animal) GetAge() int {
	return animal.age
}

func (animal *Animal) HappyBirthday() {
	animal.age++
}

func (parrot Parrot) GetAgeInHumanYears() int {
	return parrot.age * 5
}

func (parrot Parrot) Express() {
	fmt.Printf(`%s: "Hello!"`, parrot.GetName())
	fmt.Println(`
   \
   (o>
\\_//)
 \_/_)
  _|_
	`)
}

func (turtle Turtle) GetAgeInHumanYears() int {
	if turtle.age <= 0 {
		return 0
	}
	return 15 + (turtle.age-1)*2
}

func (turtle Turtle) Express() {
	fmt.Println(`
  _____     ____
 /      \  |  o |
|        |/ ___\|
|_________/
|_|_| |_|_|
	`)
}

func (person *Person) SetPet(pet Pet) {
	person.pet = pet
}

func (person *Person) GetPetInfo() {
	fmt.Println("--------------------")
	if person.pet == nil {
		fmt.Printf("%s has no pet yet\n", person.name)
		fmt.Println("--------------------")
		return
	}

	fmt.Printf("Pet for %s is named %s\n", person.name, person.pet.GetName())
	fmt.Printf(
		"Age: %d year(s), about %d human years\n",
		person.pet.GetAge(),
		person.pet.GetAgeInHumanYears(),
	)
	fmt.Println("And it expresses itself like this:")
	person.pet.Express()
	fmt.Println("--------------------")
}

func main() {
	kesha := Parrot{
		Animal: Animal{
			name: "Kesha",
			age:  2,
		},
	}

	zavala := Turtle{
		Animal: Animal{
			name: "Zavala",
			age:  4,
		},
	}

	kesha.HappyBirthday()

	Jack := Person{
		name: "Jack",
		pet:  &kesha,
	}

	John := Person{
		name: "John",
	}

	Jack.GetPetInfo()
	John.GetPetInfo()

	John.SetPet(&zavala)

	John.GetPetInfo()
}
