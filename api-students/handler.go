package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
)

type StudentHandler struct {
	repo repository.StudentRepository
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func terjemahkanError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "Student tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "NIM sudah digunakan")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

func (h *StudentHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	students, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Gagal mengambil data student")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return okList(c, "Daftar Student berhasil diambil", students, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *StudentHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	student, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Gagal mengambil data student")
	}

	return ok(c, "Student ditemukan", student)
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

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
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru, err := h.repo.Create(ctx, model.Student{
		Name:     req.Name,
		NIM:      req.NIM,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "Gagal menyimpan student")
	}

	return created(c, "Student berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.NIM = strings.TrimSpace(req.NIM)

	if req.Name == "" {
		errs["name"] = "Wajib diisi pada PUT"
	}
	if req.NIM == "" {
		errs["nim"] = "Wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "Harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.Student{
		ID: id, Name: req.Name, NIM: req.NIM, Grade: req.Grade, IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui student")
	}

	return ok(c, "Student berhasil diganti seluruhnya", hasil)
}

func (h *StudentHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	if req.Name == nil && req.NIM == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "Tidak ada field yang diubah")
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui student")
	}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return failValidation(c, map[string]string{"name": "Tidak boleh kosong"})
		}
		saatIni.Name = *req.Name
	}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			return failValidation(c, map[string]string{"nim": "Tidak boleh kosong"})
		}
		saatIni.NIM = *req.NIM
	}

	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			return failValidation(c, map[string]string{"grade": "Harus antara 0 dan 100"})
		}
		saatIni.Grade = *req.Grade
	}

	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}

	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui student")
	}

	return ok(c, "Student berhasil diperbarui sebagian", hasil)
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "Gagal menghapus student")
	}

	return noContent(c)
}
