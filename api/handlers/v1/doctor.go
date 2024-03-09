package v1

import (
	"context"
	"fmt"
	"myproject/admin-api-gateway/api/handlers/tokens"
	"myproject/admin-api-gateway/api/models"
	"myproject/admin-api-gateway/pkg/etc"
	"net/http"
	"time"

	pb "myproject/admin-api-gateway/genproto/healthcare-service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateDoctor
// @Router /v1/doctor/create [post]
// @Security BearerAuth
// @Summary create doctor
// @Tags Doctor
// @Description Create a new doctor with the provided details
// @Accept json
// @Produce json
// @Param DoctorInfo body models.DoctorReq true "Create doctor"
// @Success 201 {object} models.DoctorModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateDoctor(c *gin.Context) {
	var (
		body       models.DoctorReq
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	body.Password, err = etc.GenerateHashPassword(body.Password)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while hashing password") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	body.ID = uuid.New().String()
	h.jwtHandler = tokens.JWTHandler{
		Sub:       body.ID,
		Role:      "doctor",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		TimeOut:   h.cfg.AccessTokenTimeOut,
	}
	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while generating access and refresh token") {
		return
	}

	createReq := &pb.Doctor{
		Id:            body.ID,
		FullName:      body.LastName + " " + body.FirstName,
		BirthDate:     body.BirthDate,
		Gender:        body.Gender,
		PhoneNumber:   body.PhoneNumber,
		Email:         body.Email,
		Address:       body.Address,
		Salary:        float32(body.Salary),
		Biography:     body.Biography,
		StartWorkYear: body.StartWorkYear,
		EndWorkYear:   body.EndWorkYear,
		WorkYears:     body.WorkYears,
		DepartmentId:  body.DepartmentId,
		SpecIds:       body.SpecIds,
		RefreshToken:  refresh,
	}

	respDoctor, err := h.serviceManager.HealthCareService().CreateDoctor(ctx, createReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while creating a doctor") {
		return
	}

	response := models.DoctorModel{
		ID:            respDoctor.Id,
		FullName:      respDoctor.FullName,
		BirthDate:     respDoctor.BirthDate,
		Gender:        respDoctor.Gender,
		PhoneNumber:   respDoctor.PhoneNumber,
		Email:         respDoctor.Email,
		Address:       respDoctor.Address,
		Salary:        float64(respDoctor.Salary),
		Biography:     respDoctor.Biography,
		StartWorkYear: respDoctor.StartWorkYear,
		EndWorkYear:   respDoctor.EndWorkYear,
		WorkYears:     respDoctor.WorkYears,
		DepartmentId:  respDoctor.DepartmentId,
		SpecIds:       respDoctor.SpecIds,
		AccessToken:   access,
		CreatedAt:     respDoctor.CreatedAt,
		UpdatedAt:     respDoctor.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// Get Doctor By Id
// @Router /v1/doctor/{id} [get]
// @Security BearerAuth
// @Summary get doctor by id
// @Tags Doctor
// @Description Get user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.DoctorResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetDoctorById(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	respDoctor, err := h.serviceManager.HealthCareService().GetDoctorById(ctx, &pb.GetReqStr{Id: id})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while getting doctor by id") {
		return
	}

	response := models.DoctorResp{
		ID:            respDoctor.Id,
		FullName:      respDoctor.FullName,
		BirthDate:     respDoctor.BirthDate,
		Gender:        respDoctor.Gender,
		PhoneNumber:   respDoctor.PhoneNumber,
		Email:         respDoctor.Email,
		Password:      respDoctor.Password,
		Address:       respDoctor.Address,
		Salary:        float64(respDoctor.Salary),
		Biography:     respDoctor.Biography,
		StartWorkYear: respDoctor.StartWorkYear,
		EndWorkYear:   respDoctor.EndWorkYear,
		WorkYears:     respDoctor.WorkYears,
		DepartmentId:  respDoctor.DepartmentId,
		SpecIds:       respDoctor.SpecIds,
		CreatedAt:     respDoctor.CreatedAt,
		UpdatedAt:     respDoctor.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Update doctor
// @Router /v1/doctor/update/{id} [put]
// @Security BearerAuth
// @Summary update doctor
// @Tags Doctor
// @Description Update doctor
// @Accept json
// @Produce json
// @Param id path string false "id"
// @Param UserInfo body models.DoctorUpdateReq true "Update Doctor"
// @Success 201 {object} models.DoctorResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateDoctor(c *gin.Context) {
	var (
		body        models.DoctorUpdateReq
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	updateReq := &pb.Doctor{
		Id:            body.ID,
		FullName:      body.FullName,
		BirthDate:     body.BirthDate,
		Address:       body.Address,
		Salary:        float32(body.Salary),
		Biography:     body.Biography,
		StartWorkYear: body.StartWorkYear,
		EndWorkYear:   body.EndWorkYear,
		WorkYears:     body.WorkYears,
		DepartmentId:  body.DepartmentId,
	}

	if updateReq.Id == "" {
		id := c.Param("id")
		if id == "" {
			if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("id is required"), ErrorBadRequest) {
				return
			}
		}
		updateReq.Id = id
	}

	respDoctor, err := h.serviceManager.HealthCareService().UpdateDoctor(ctx, updateReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	response := models.DoctorResp{
		ID:            respDoctor.Id,
		FullName:      respDoctor.FullName,
		BirthDate:     respDoctor.BirthDate,
		Gender:        respDoctor.Gender,
		PhoneNumber:   respDoctor.PhoneNumber,
		Email:         respDoctor.Email,
		Password:      respDoctor.Password,
		Address:       respDoctor.Address,
		Salary:        float64(respDoctor.Salary),
		Biography:     respDoctor.Biography,
		StartWorkYear: respDoctor.StartWorkYear,
		EndWorkYear:   respDoctor.EndWorkYear,
		WorkYears:     respDoctor.WorkYears,
		DepartmentId:  respDoctor.DepartmentId,
		SpecIds:       respDoctor.SpecIds,
		CreatedAt:     respDoctor.CreatedAt,
		UpdatedAt:     respDoctor.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Delete Doctor
// @Router /v1/doctor/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete doctor
// @Tags Doctor
// @Description Delete doctor
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteDoctor(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err := h.serviceManager.HealthCareService().DeleteDoctor(ctx, &pb.GetReqStr{Id: id})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while deleting doctor") {
		return
	}

	c.JSON(http.StatusOK, models.Status{Message: "doctor was successfully deleted"})
}

// List doctors
// @Router /v1/doctors/{page}/{limit} [get]
// @Security BearerAuth
// @Summary get doctors' list
// @Tags Doctor
// @Description get all doctors
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Success 201 {object} models.ListDoctors
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListDoctors(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page, err := ParsePageQueryParam(c)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}
	limit, err := ParseLimitQueryParam(c)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.HealthCareService().GetAllDoctors(ctx, &pb.GetAll{Page: int64(page), Limit: int64(limit)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}
