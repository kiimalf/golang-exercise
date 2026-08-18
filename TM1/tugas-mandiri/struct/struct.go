package main

import "fmt"

type Student struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

func (s Student) GetInfo() string {
	return fmt.Sprintf("name: %s, grade: %.2f, status: %v", s.Name, s.Grade, s.IsActive)
}

func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

func (s *Student) Activate() {
	s.IsActive = true
}

func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	s := Student{ID: 1, Name: "Hakim"}
	fmt.Println(s.GetInfo())

	s.UpdateGrade(9.45)
	fmt.Printf("\nSetelah UpdateGrade 9.45 :\n")
	fmt.Println(s.GetInfo())

	s.Activate()
	fmt.Printf("\nSetelah Activate :\n")
	fmt.Println(s.GetInfo())

	s.Deactivate()
	fmt.Printf("\nSetelah Deactivate :\n")
	fmt.Println(s.GetInfo())
}
