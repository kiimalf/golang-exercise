package service

import (
	"api-students/app/model"
	"testing"
)

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
	}

	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{
		ID:       1,
		Name:     "kunto",
		NIM:      "99999",
		Grade:    75.3,
		IsActive: true,
	}
	inactive := false

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{
		IsActive: &inactive,
	})

	if len(errs) != 0 {
		t.Fatalf("Tidak seharusnya ada error: %v", errs)
	}
	if result.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}
	if result.Name != "kunto" {
		t.Error("Field yang tidak dikirim seharusnya tidak berubah")
	}
}
