package service

import (
	"strings"

	"api-students/app/model"
)

func ValidateCreate(req model.CreateStudentRequest) map[string]string {
	errs := map[string]string{}

	req.Name = strings.TrimSpace(req.Name)
	req.NIM = strings.TrimSpace(req.NIM)

	if req.Name == "" {
		errs["name"] = "Wajib diisi"
	}
	if req.NIM == "" {
		errs["nim"] = "Wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "Harus antara 0 dan 100"
	}

	return errs
}

func ValidateReplace(req model.ReplaceStudentRequest) map[string]string {
	errs := map[string]string{}

	if req.Name == "" {
		errs["name"] = "Wajib diisi pada PUT"
	}
	if req.NIM == "" {
		errs["nim"] = "Wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "Harus antara 0 dan 100"
	}

	return errs
}

func ApplyPatch(
	current model.Student, req model.PatchStudentRequest,
) (model.Student, map[string]string) {
	errs := map[string]string{}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs["name"] = "Tidak boleh kosong"
		}
		current.Name = *req.Name
	}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			errs["nim"] = "Tidak boleh kosong"
		}
		current.NIM = *req.NIM
	}

	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "Harus antara 0 dan 100"
		}
		current.Grade = *req.Grade
	}

	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}

func IsEmptyPatch(req model.PatchStudentRequest) bool {
	return req.Name == nil && req.NIM == nil && req.Grade == nil && req.IsActive == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}

	return (total + limit - 1) / limit
}
