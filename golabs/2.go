package main

import (
	"errors"
	"fmt"
)

// 1. Функция hello, которая принимает строку и возвращает "What's up, name?"
func hello(name string) string {
	return fmt.Sprintf("What's up, %s!", name)
}

// 2. Функция printEven, которая выводит чётные числа в заданном диапазоне
func printEven(a, b int64) error {
	if a > b {
		return errors.New("левая граница диапазона больше правой")
	}

	for i := a; i <= b; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}
	return nil
}

// 3. Функция apply, которая выполняет действия с двумя числами
func apply(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, errors.New("действие не поддерживается")
	}
}

func main() {
	// Тестирование функции hello
	fmt.Println("=== Тест функции hello ===")
	fmt.Println(hello("Christopher"))

	// Тестирование функции printEven
	fmt.Println("=== Тест функции printEven ===")
	if err := printEven(2, 10); err != nil {
		fmt.Println("Ошибка:", err)
	}

	if err := printEven(10, 2); err != nil {
		fmt.Println("Ошибка:", err)
	}

	// Тестирование функции apply
	fmt.Println("=== Тест функции apply ===")
	result, err := apply(3, 5, "+")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("3 + 5 = %f\n", result)
	}

	result, err = apply(7, 10, "*")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("7 * 10 = %f\n", result)
	}

	result, err = apply(3, 5, "#")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("3 # 5 = %f\n", result)
	}

	result, err = apply(10, 0, "/")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("10 / 0 = %f\n", result)
	}
}
