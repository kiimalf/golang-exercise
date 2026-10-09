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
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	NIM      string  `json:"nim" validate:"required,nim"`
	Grade    float64 `json:"grade" validate:"gte=0,lte=100"`
	IsActive *bool   `json:"is_active"`
}

type ReplaceStudentRequest struct {
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	NIM      string  `json:"nim" validate:"required,nim"`
	Grade    float64 `json:"grade" validate:"gte=0,lte=100"`
	IsActive bool    `json:"is_active"`
}

type PatchStudentRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitnil,min=3,max=100"`
	NIM      *string  `json:"nim,omitempty" validate:"omitnil,nim"`
	Grade    *float64 `json:"grade,omitempty" validate:"omitnil,gte=0,lte=100"`
	IsActive *bool    `json:"is_active,omitempty"`
}
