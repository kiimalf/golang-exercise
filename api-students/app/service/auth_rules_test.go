package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestValidateRegister_PasswordTooShort(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "abc",
	}
	errs := helper.ValidateStruct(req)
	if errs["password"] != "Minimal 8 karakter" {
		t.Errorf("Harap error 'minimal 8 karakter', dapat %q", errs["password"])
	}
}

func TestValidateRegister_PasswordNoDigit(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "passwordrahasia",
	}
	errs := helper.ValidateStruct(req)
	if errs["password"] != "Harus memuat huruf dan angka" {
		t.Errorf("Harap error 'harus memuat huruf dan angka', dapat %q", errs["password"])
	}
}

func TestValidateRegister_PasswordValid(t *testing.T) {
	req := model.RegisterRequest{
		Username: "budi",
		Email:    "budi@example.com",
		Password: "rahasia123",
	}
	errs := helper.ValidateStruct(req)
	if _, exists := errs["password"]; exists {
		t.Errorf("Password valid seharusnya tidak menghasilkan error password, dapat %q", errs["password"])
	}
}
