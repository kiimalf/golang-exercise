package service

import (
	"strings"

	"api-students/app/model"
)

func ValidateCreateUser(req model.CreateUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "Wajib diisi"
	}
	if !isValidEmail(req.Email) {
		errs["email"] = "Format email tidak valid"
	}
	if len(req.Password) < 8 {
		errs["password"] = "Minimal 8 karakter"
	}

	return errs
}

func ValidateReplaceUser(req model.ReplaceUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "Wajib diisi pada PUT"
	}
	if !isValidEmail(req.Email) {
		errs["email"] = "Wajib diisi dan berformat email pada PUT"
	}

	return errs
}

func ApplyPatchUser(current model.User, req model.PatchUserRequest) (model.User, map[string]string) {
	errs := map[string]string{}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			errs["username"] = "Tidak boleh kosong"
		} else {
			current.Username = *req.Username
		}
	}

	if req.Email != nil {
		if !isValidEmail(*req.Email) {
			errs["email"] = "Format email tidak valid"
		} else {
			current.Email = *req.Email
		}
	}

	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}

func IsEmptyPatchUser(req model.PatchUserRequest) bool {
	return req.Username == nil && req.Email == nil && req.IsActive == nil
}
