package service

import (
	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SupervisorService struct {
	supervisors repository.SupervisorRepository
	students    repository.StudentRepository
	tx          *repository.TxManager
}

func NewSupervisorService(
	supervisors repository.SupervisorRepository,
	students repository.StudentRepository,
	tx *repository.TxManager,
) *SupervisorService {
	return &SupervisorService{
		supervisors: supervisors,
		students:    students,
		tx:          tx,
	}
}

func (s *SupervisorService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	list, err := s.supervisors.FindAll(ctx)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, fiber.StatusOK, "Daftar pembimbing berhasil diambil", list)
}

func (s *SupervisorService) Assign(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	supervisorID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID harus berupa angka positif")
	}

	var req model.AssignSupervisorRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	var updated model.Supervisor
	err := s.tx.WithTx(ctx, func(tx repository.DBTX) error {
		supervisorRepo := s.supervisors.WithTx(tx)
		studentRepo := s.students.WithTx(tx)

		if err := studentRepo.LockStudent(ctx, req.StudentID); err != nil {
			return err
		}

		var err error
		updated, err = supervisorRepo.IncrementSupervised(ctx, supervisorID)
		if err != nil {
			return err
		}

		return supervisorRepo.InsertSupervision(ctx, req.StudentID, supervisorID)
	})

	if err != nil {
		return translateSupervisorError(err)
	}
	return helper.Created(c, "Berhasil menetapkan pembimbing", updated,
		"/api/v1/supervisors/"+strconv.Itoa(supervisorID))
}

func (s *SupervisorService) Switch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.SwitchSupervisorRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}
	if req.FromSupervisorID == req.ToSupervisorID {
		return helper.Validation(map[string]string{
			"to_supervisor_id": "Tidak boleh sama dengan from_supervisor_id",
		})
	}

	var joined model.Supervisor
	err := s.tx.WithTx(ctx, func(tx repository.DBTX) error {
		supervisorRepo := s.supervisors.WithTx(tx)
		studentRepo := s.students.WithTx(tx)

		if err := studentRepo.LockStudent(ctx, req.StudentID); err != nil {
			return err
		}

		if err := supervisorRepo.LockSupervisors(ctx, req.FromSupervisorID, req.ToSupervisorID); err != nil {
			return err
		}

		if err := supervisorRepo.DeleteSupervision(ctx, req.StudentID, req.FromSupervisorID); err != nil {
			return err
		}
		if _, err := supervisorRepo.DecrementSupervised(ctx, req.FromSupervisorID); err != nil {
			return err
		}

		var err error
		joined, err = supervisorRepo.IncrementSupervised(ctx, req.ToSupervisorID)
		if err != nil {
			return err
		}
		return supervisorRepo.InsertSupervision(ctx, req.StudentID, req.ToSupervisorID)
	})

	if err != nil {
		return translateSupervisorError(err)
	}
	return helper.Success(c, fiber.StatusOK, "Berhasil mengganti pembimbing", joined)
}

func (s *SupervisorService) UpdateQuota(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID harus berupa angka positif")
	}

	var req model.UpdateSupervisorQuotaRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	updated, err := s.supervisors.UpdateQuota(ctx, id, req.Quota, req.Version)
	if err != nil {
		return translateSupervisorError(err)
	}
	return helper.Success(c, fiber.StatusOK, "Kuota pembimbing berhasil diubah", updated)
}

func translateSupervisorError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("data tidak ditemukan")
	case errors.Is(err, repository.ErrSupervisorFull):
		return helper.NewError(fiber.StatusConflict, "SUPERVISOR_FULL", "kuota dosen pembimbing sudah penuh")
	case errors.Is(err, repository.ErrAlreadySupervised):
		return helper.NewError(fiber.StatusConflict, "ALREADY_SUPERVISED", "mahasiswa sudah memiliki pembimbing")
	case errors.Is(err, repository.ErrNotSupervised):
		return helper.NewError(fiber.StatusConflict, "NOT_SUPERVISED", "mahasiswa tidak dibimbing oleh dosen tersebut")
	case errors.Is(err, repository.ErrVersionConflict):
		return helper.NewError(fiber.StatusConflict, "VERSION_CONFLICT", "data sudah diubah pihak lain, muat ulang lalu coba lagi")
	case errors.Is(err, repository.ErrQuotaBelowFilled):
		return helper.NewError(fiber.StatusConflict, "QUOTA_BELOW_SUPERVISED", "kuota tidak boleh lebih kecil dari mahasiswa yang sedang dibimbing")
	default:
		return helper.Internal(err)
	}
}
