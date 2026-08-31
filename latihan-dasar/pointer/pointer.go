package main

import (
	"fmt"
)

func swap(a, b *int) {
	*a, *b = *b, *a
}

func swapByValue(c, d int) {
	c, d = d, c
	fmt.Printf("didalam function swapByValue, c: %v, d: %v\n", c, d)
}

func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

func main() {
	a, b := 1, 2
	fmt.Printf("sebelum swap, a: %v, b: %v\n", a, b)

	swap(&a, &b)
	fmt.Printf("setelah swap, a: %v, b: %v\n\n", a, b)

	c, d := 3, 4
	swap(&a, &b)
	fmt.Printf("sebelum swapByValue, c: %v, d: %v\n", c, d)

	swapByValue(c, d)
	fmt.Printf("setelah swapByValue, c: %v, d: %v\n\n", c, d)

	buah := []string{"Apel", "Mangga"}
	fmt.Println("sebelum updateSlice: ", buah)
	updateSlice(&buah, "Anggur")
	fmt.Println("Setelah updateSlice: ", buah)
}
