package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type UserService struct {
	repo  repository.UserRepository
	perms *helper.PermissionSet
}

func NewUserService(
	repo repository.UserRepository,
	perms *helper.PermissionSet,
) *UserService {
	return &UserService{repo: repo, perms: perms}
}

func (s *UserService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	format, err := helper.Negotiate(c, helper.FormatJSON, helper.FormatCSV)
	if err != nil {
		return err
	}

	q, err := helper.ParseCursorQuery(c)
	if err != nil {
		return err
	}

	rows, err := s.repo.FindAfterCursor(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}

	hasMore := len(rows) > q.Limit
	if hasMore {
		rows = rows[:q.Limit]
	}

	if format == helper.FormatCSV {
		return helper.WriteUsersCSV(c, rows)
	}

	meta := &model.CursorMeta{Limit: q.Limit, HasMore: hasMore}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		meta.NextCursor = helper.EncodeCursor(last.CreatedAt, last.ID)
	}

	return helper.SuccessCursor(c, "Daftar user berhasil diambil", rows, meta)
}

func (s *UserService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("Id harus berupa angka positif")
	}

	if !CanAccessUser(current, id, s.perms, "user:read:any") {
		return helper.Forbidden("Tidak berhak mengakses data user lain")
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateUserError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "User ditemukan", user)
}

func (s *UserService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	newUser, err := s.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	})

	if err != nil {
		return translateUserError(err, "user")
	}

	return helper.Created(c, "User berhasil dibuat", newUser,
		"/api/v1/users/"+strconv.Itoa(newUser.ID))
}

func (s *UserService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("Id harus berupa angka positif")
	}

	if !CanAccessUser(current, id, s.perms, "user:update:any") {
		return helper.Forbidden("Tidak berhak mengakses data user lain")
	}

	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	result, err := s.repo.Update(ctx, model.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		IsActive: req.IsActive,
	})

	if err != nil {
		return translateUserError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "User berhasil diganti seluruhnya", result)
}

func (s *UserService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("Id harus berupa angka positif")
	}

	if !CanAccessUser(currentUser, id, s.perms, "user:update:any") {
		return helper.Forbidden("Tidak berhak mengakses data user lain")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}

	if IsEmptyPatchUser(req) {
		return helper.BadRequest("Tidak ada field yang diubah")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	targetUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateUserError(err, "user")
	}

	updated := ApplyPatchUser(targetUser, req)

	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		return translateUserError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "User berhasil diperbarui sebagian", result)
}

func (s *UserService) AssignRole(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("Id harus berupa angka positif")
	}

	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}

	if errs := ValidateAssignRole(currentUser, id, req, s.perms); len(errs) > 0 {
		return helper.Validation(errs)
	}

	result, err := s.repo.UpdateRole(ctx, id, strings.TrimSpace(req.Role))
	if err != nil {
		return translateUserError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "Role user berhasil diubah", result)
}

func (s *UserService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("Id harus berupa angka positif")
	}

	if current.UserID == id {
		return helper.Forbidden("Tidak boleh menghapus akun sendiri")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateUserError(err, "user")
	}

	return helper.NoContent(c)
}

func translateUserError(err error, entity string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound(entity + " tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Conflict("Username sudah dipakai")
	default:
		return helper.Internal(err)
	}
}
