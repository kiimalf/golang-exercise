package model

import "time"

type Student struct {
	ID        int       `json:"id"`
	OwnerID   *int      `json:"owner_id"`
	Name      string    `json:"name"`
	NIM       string    `json:"nim"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStudentRequest struct {
	Name     string  `json:"name"`
	NIM      string  `json:"nim"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

type ReplaceStudentRequest struct {
	Name     string  `json:"name"`
	NIM      string  `json:"nim"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

type PatchStudentRequest struct {
	Name     *string  `json:"name,omitempty"`
	NIM      *string  `json:"nim,omitempty"`
	Grade    *float64 `json:"grade,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}
