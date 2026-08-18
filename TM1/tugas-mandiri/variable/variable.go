package main

import (
	"fmt"
)

func main() {
	//variable
	nama := "Nabil Hakim Alfikri"
	fmt.Printf("string nama: %v\n", nama)

	nim := 434241055
	fmt.Printf("int nim: %v\n", nim)

	var ipk float64 = 3.85
	fmt.Printf("float64 ipk: %.2f\n", ipk)

	isActive := true
	fmt.Printf("bool isActive: %v\n", isActive)

	mataKuliah := []string{"Machine Learning", "Design Thinking"}
	fmt.Println("Daftar Matakuliah: ", mataKuliah)

	//declare map
	nilaiMahasiswa := make(map[string]int)

	//menambah data map
	nilaiMahasiswa["Anto"] = 90
	nilaiMahasiswa["Anti"] = 85
	nilaiMahasiswa["Anta"] = 60

	//pengecekan keberadaan nilai dalam map
	namaDicari := "Anto"
	if nilai, ada := nilaiMahasiswa[namaDicari]; ada {
		fmt.Printf("Nilai %v: %v\n", namaDicari, nilai)
	} else {
		fmt.Printf("Data %v tidak ditemukan\n", namaDicari)
	}

	namaDicari2 := "Dina"
	if nilai, ada := nilaiMahasiswa[namaDicari2]; ada {
		fmt.Printf("Nilai %v: %v\n", namaDicari2, nilai)
	} else {
		fmt.Printf("Data %v tidak ditemukan\n", namaDicari2)
	}

	//menghapus nilai dalam map
	delete(nilaiMahasiswa, "Anti")

	//Menampilkan seluru isi map
	for key, value := range nilaiMahasiswa {
		fmt.Printf("%v: %v\n", key, value)
	}

}
