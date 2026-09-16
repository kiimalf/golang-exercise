package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateRegister_PasswordTooShort(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "abc",
	}
	errs := ValidateRegister(req)
	if errs["password"] != "minimal 8 karakter" {
		t.Errorf("Harap error 'minimal 8 karakter', dapat %q", errs["password"])
	}
}

func TestValidateRegister_PasswordNoDigit(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "passwordrahasia",
	}
	errs := ValidateRegister(req)
	if errs["password"] != "harus memuat huruf dan angka" {
		t.Errorf("Harap error 'harus memuat huruf dan angka', dapat %q", errs["password"])
	}
}

func TestValidateRegister_PasswordValid(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "rahasia123",
	}
	errs := ValidateRegister(req)
	if _, exists := errs["password"]; exists {
		t.Errorf("Password valid seharusnya tidak menghasilkan error password, dapat %q", errs["password"])
	}
}
