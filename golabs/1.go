package main

import (
	"errors"
	"fmt"
	"math"
)

// ----------------------- Задание 1: Массивы и срезы -----------------------

// formatIP принимает массив из 4 байтов и возвращает строку с IP-адресом
func formatIP(ip [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// listEven возвращает срез четных чисел в диапазоне и ошибку
func listEven(start, end int) ([]int, error) {
	if start > end {
		return nil, errors.New("левая граница больше правой")
	}

	var evens []int
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			evens = append(evens, i)
		}
	}
	return evens, nil
}

// ----------------------- Задание 2: Карты -----------------------

// countCharacters принимает строку и подсчитывает количество вхождений каждого символа
func countCharacters(s string) map[rune]int {
	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
	}
	return charCount
}

// ----------------------- Задание 3: Структуры, методы и интерфейсы -----------------------

// Point определяет структуру для точки на плоскости с координатами X и Y
type Point struct {
	X, Y float64
}

// LineSegment определяет структуру для отрезка, состоящего из двух точек: начала и конца
type LineSegment struct {
	Start, End Point
}

// Length метод для вычисления длины отрезка
func (ls LineSegment) Length() float64 {
	// Используем формулу расстояния между двумя точками: sqrt((x2 - x1)^2 + (y2 - y1)^2)
	dx := ls.End.X - ls.Start.X
	dy := ls.End.Y - ls.Start.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// ----------------------- Задание 3: Треугольник и Круг -----------------------

// Triangle определяет структуру треугольника, используя три точки
type Triangle struct {
	A, B, C Point
}

// Circle определяет структуру круга с центром и радиусом
type Circle struct {
	Center Point
	Radius float64
}

// Метод Area для треугольника возвращает площадь треугольника
func (t Triangle) Area() float64 {
	// Используем формулу Герона для вычисления площади треугольника
	a := distance(t.A, t.B)
	b := distance(t.B, t.C)
	c := distance(t.C, t.A)
	s := (a + b + c) / 2
	return math.Sqrt(s * (s - a) * (s - b) * (s - c))
}

// Метод Area для круга возвращает площадь круга
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Функция для вычисления расстояния между двумя точками
func distance(p1, p2 Point) float64 {
	return math.Sqrt((p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y))
}

// Shape интерфейс с методом Area
type Shape interface {
	Area() float64
}

// Функция для вывода площади фигуры, принимающая интерфейс Shape и название фигуры
func printArea(s Shape, figureName string) {
	result := s.Area()
	fmt.Printf("Площадь %s: %.2f\n", figureName, result)
}

func main() {
	// ----------------------- Задание 1: Массивы и срезы -----------------------

	// Пример использования функции formatIP
	ip := [4]byte{192, 168, 0, 1}
	fmt.Println("IP-адрес:", formatIP(ip)) // вывод: 192.168.0.1

	// Пример использования функции listEven
	evens, err := listEven(3, 10)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Четные числа в диапазоне 3-10:", evens) // вывод: [4 6 8 10]
	}

	// Пример с ошибкой
	evens, err = listEven(10, 3)
	if err != nil {
		fmt.Println("Ошибка:", err) // вывод: Ошибка: левая граница больше правой
	} else {
		fmt.Println("Четные числа:", evens)
	}

	// ----------------------- Задание 2: Карты -----------------------

	// Пример использования функции countCharacters
	str := "hello world"
	charCount := countCharacters(str)
	fmt.Println("Подсчёт символов в строке 'hello world':")
	for char, count := range charCount {
		fmt.Printf("'%c' встречается %d раз(а)\n", char, count)
	}

	// ----------------------- Задание 3: Структуры, методы и интерфейсы -----------------------

	// Пример использования структуры Point
	startPoint := Point{X: 0, Y: 0}
	endPoint := Point{X: 3, Y: 4}

	// Пример использования структуры LineSegment
	line := LineSegment{Start: startPoint, End: endPoint}

	// Вычисляем длину отрезка
	length := line.Length()
	fmt.Printf("Длина отрезка с началом в (%.1f, %.1f) и концом в (%.1f, %.1f) равна %.2f\n",
		line.Start.X, line.Start.Y, line.End.X, line.End.Y, length) // вывод: 5.00

	// ----------------------- Задание 3: Треугольник и Круг -----------------------

	// Пример использования структуры Triangle (треугольник)
	triangle := Triangle{
		A: Point{X: 0, Y: 0},
		B: Point{X: 3, Y: 0},
		C: Point{X: 0, Y: 4},
	}

	// Пример использования структуры Circle (круг)
	circle := Circle{
		Center: Point{X: 0, Y: 0},
		Radius: 5,
	}

	// Выводим площадь треугольника
	printArea(triangle, "треугольника") // вывод: Площадь треугольника: 6.00

	// Выводим площадь круга
	printArea(circle, "круга") // вывод: Площадь круга: 78.54

	// ----------------------- Задание 4: Функциональное программирование -----------------------

	// Функция Map принимает срез и функцию, применяет функцию ко всем элементам среза
	Map := func(slice []float64, fn func(float64) float64) []float64 {
		result := make([]float64, len(slice))
		for i, v := range slice {
			result[i] = fn(v) // Применяем функцию к каждому элементу
		}
		return result
	}

	// Функция для возведения числа в квадрат
	square := func(x float64) float64 {
		return x * x
	}

	// Создаем срез значений
	values := []float64{1, 2, 3, 4, 5}

	// Выводим исходный срез
	fmt.Println("Исходный срез:", values)

	// Применяем функцию Map с функцией возведения в квадрат
	squaredValues := Map(values, square)

	// Выводим срез после применения функции
	fmt.Println("Срез после возведения в квадрат:", squaredValues)
}
