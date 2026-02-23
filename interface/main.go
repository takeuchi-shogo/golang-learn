package main

import (
	"fmt"
	"log/slog"
)

type MyType interface {
	~string | ~int | ~float64 | ~bool
}

func PrintMyType[T MyType](t T) {
	switch v := any(t).(type) {
	case string:
		fmt.Println("string: ", v)
	case int:
		fmt.Println("int: ", v)
	case float64:
		fmt.Println("float64: ", v)
	case bool:
		fmt.Println("bool: ", v)
	default:
		fmt.Println("unknown type")
	}
}

type Dog struct {
	name string
}

type Cat struct {
	name string
}

type Bird struct {
	name string
}

type Animal interface {
	// MakeSound は、動物の鳴き声を返すメソッド
	// 返り値は string で鳴き声を返す
	MakeSound() string
	// GetName は、動物の名前を返すメソッド
	// 返り値は string で動物の名前を返す
	GetName() string
}

type Organism interface {
	Animal
}

type organism struct {
	name string
	Animal
}

func (o *organism) GetName() string {
	return o.name
}

func NewOrganism(name string, animal Animal) Organism {
	return &organism{name: name, Animal: animal}
}

// コンパイル時にインターフェース実装をチェック（マジック的なコード）
var _ Animal = (*Dog)(nil)
var _ Animal = (*Cat)(nil)
var _ Animal = (*Bird)(nil)

func (d *Dog) MakeSound() string {
	return "Woof"
}

func (d *Dog) GetName() string {
	return d.name
}

func (c *Cat) MakeSound() string {
	return "Meow"
}

func (c *Cat) GetName() string {
	return c.name
}

func (b *Bird) MakeSound() string {
	return "Tweet"
}

func (b *Bird) GetName() string {
	return b.name
}

func PrintAnimalSound(animal Animal) {
	fmt.Println(animal.MakeSound())
	fmt.Println(animal.GetName())
}

type AppError interface {
	// Error() string
}

const (
	ErrorCodeInvalidRequest = "invalid_request"
)

type appErrorImpl struct {
	code       string
	message    string
	statusCode int
}

func (e *appErrorImpl) Error() string {
	return fmt.Sprintf("code: %s, message: %s, statusCode: %d", e.code, e.message, e.statusCode)
}

func sample() error {
	return &appErrorImpl{message: "App error", code: ErrorCodeInvalidRequest, statusCode: 400}
}

type Email string

func (e Email) String() string {
	return "xxxxxxx@xxxxxxx.com"
}

func main() {
	dog := &Dog{name: "Buddy"}
	cat := &Cat{name: "Whiskers"}
	bird := &Bird{name: "Tweetie"}
	PrintAnimalSound(dog)
	PrintAnimalSound(cat)
	PrintAnimalSound(bird)

	organism := NewOrganism("Organism", dog)
	PrintAnimalSound(organism)

	err := sample()
	if err != nil {
		fmt.Println(err.Error())
	}

	email := Email("test@test.com")
	fmt.Println(email.String())
	slog.Info("email", slog.String("email", email.String()))
}
