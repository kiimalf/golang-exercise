package service

import (
	"strings"

	"latihan-fiber/app/model"
)

func ApplyPatch(current model.User, req model.PatchUserRequest) model.User {
	if req.Username != nil {
		current.Username = strings.TrimSpace(*req.Username)
	}

	if req.Email != nil {
		current.Email = strings.TrimSpace(*req.Email)
	}

	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current
}

func IsEmptyPatch(req model.PatchUserRequest) bool {
	return req.Username == nil && req.Email == nil && req.IsActive == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}

	return (total + limit - 1) / limit
}
