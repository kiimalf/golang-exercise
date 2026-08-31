package main

import (
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (u *User) Activate()   { u.IsActive = true }
func (u User) Info() string { return fmt.Sprintf("%s aktif: %v", u.Username, u.IsActive) }

func tukarNilai(a, b *int) { *a, *b = *b, *a }

func main() {
	//variable
	var nama string = "Sari"
	umur := 20
	var kosong string
	fmt.Printf("%s %d zero value: %q\n", nama, umur, kosong)

	//slice
	nilai := []int{10, 20, 30}
	nilai = append(nilai, 40)
	fmt.Println(nilai, nilai[1:3], len(nilai), cap(nilai))

	//map
	skor := make(map[string]int)
	skor["Sari"] = 90
	if n, ada := skor["Budi"]; ada {
		fmt.Println("Budi: ", n)
	} else {
		fmt.Println("Budi belum punya nilai")
	}

	//channel
	ch := make(chan string)
	go func() { ch <- "halo dari goroutine" }()
	fmt.Println(<-ch)

	//pointer
	a, b := 1, 2
	tukarNilai(&a, &b)
	fmt.Println("setelah nilai ditukar: ", a, b)

	//struct
	u := User{ID: 1, Username: "Sari"}
	u.Activate()
	fmt.Println(u.Info())
}
