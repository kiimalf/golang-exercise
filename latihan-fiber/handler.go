// Bagian 1 : Penyimpanan dan bantuan
package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func terjemahkanError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "Username sudah dipakai")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	users, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Gagal mengambil data user")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return okList(c, "Daftar user berhasil diambil", users, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Gagal mengambil data user")
	}

	return ok(c, "User ditemukan", user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	errs := map[string]string{}
	if req.Username == "" {
		errs["username"] = "Wajib diisi"
	}
	if !strings.Contains(req.Email, "@") {
		errs["Email"] = "Format email tidak valid"
	}
	if len(req.Password) < 8 {
		errs["Password"] = "Minimal 8 karakter"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru, err := h.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	})
	if err != nil {
		return terjemahkanError(c, err, "Gagal menyimpan user")
	}

	return created(c, "User berhasil dibuat", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func (h *UserHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "Wajib diisi pada PUT"
	}
	if !strings.Contains(req.Email, "@") {
		errs["Email"] = "Wajib diisi pada PUT"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.User{
		ID: id, Username: req.Username, Email: req.Email, IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui user")
	}

	return ok(c, "User berhasil diganti seluruhnya", hasil)
}

func (h *UserHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "Tidak ada field yang diubah")
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui user")
	}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"username": "Tidak boleh kosong"})
		}
		saatIni.Username = *req.Username
	}

	if req.Email != nil {
		if !strings.Contains(*req.Email, "@") {
			return failValidation(c, map[string]string{"email": "Format email tidak valid"})
		}
		saatIni.Email = *req.Email
	}

	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return terjemahkanError(c, err, "Gagal memperbarui user")
	}

	return ok(c, "User berhasil diperbarui sebagian", hasil)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "Gagal menghapus user")
	}

	return noContent(c)
}
