package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"myproject/admin-api-gateway/api/handlers/tokens"
	"myproject/admin-api-gateway/api/models"
	"myproject/admin-api-gateway/email"
	"myproject/admin-api-gateway/pkg/etc"
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "myproject/admin-api-gateway/genproto/healthcare-service"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// Register doctor
// @Router /v1/doctor/register [post]
// @Summary register doctor
// @Tags Doctor
// @Description Register a new doctor with the provided details
// @Accept json
// @Produce json
// @Param DoctorInfo body models.DoctorReq true "Register doctor"
// @Success 201 {object} models.RegisterRespModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) RegisterDoctor(c *gin.Context) {
	var (
		body       models.DoctorReq
		code       string
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	body.ID = uuid.New().String()
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	err = body.Validate()
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorValidationError) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.HealthCareService().CheckUniques(ctx, &pb.CheckUniqReq{
		Field: "email",
		Value: body.Email,
	})

	if handleInternalServerErrorWithMessage(c, h.log, err, "failed to check email uniqueness") {
		return
	}

	if exists.Status {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("you've already registered before, try to log in"), ErrorBadRequest) {
			return
		}
	}

	code = etc.GenerateCode(5)
	registerDoctor := models.RegisterDoctor{
		ID:            body.ID,
		FullName:      body.LastName + " " + body.FirstName,
		BirthDate:     body.BirthDate,
		Gender:        body.Gender,
		PhoneNumber:   body.PhoneNumber,
		Email:         body.Email,
		Address:       body.Address,
		Salary:        float64(body.Salary),
		Biography:     body.Biography,
		StartWorkYear: body.StartWorkYear,
		EndWorkYear:   body.EndWorkYear,
		WorkYears:     body.WorkYears,
		DepartmentId:  body.DepartmentId,
		SpecIds:       body.SpecIds,
		IsVerified:    false,
		Code:          code,
	}

	doctorJson, err := json.Marshal(registerDoctor)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while marshaling doctor json") {
		return
	}

	timeOut := time.Second * 300

	err = h.inMemoryStorage.SetWithTTL(registerDoctor.Email, string(doctorJson), int(timeOut.Seconds()))
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while setting with ttl to redis") {
		return
	}

	message, err := email.SendVerificationCode(email.EmailPayload{
		From:     h.cfg.SendEmailFrom,
		To:       registerDoctor.Email,
		Password: h.cfg.EmailCode,
		Code:     registerDoctor.Code,
		Message:  fmt.Sprintf("Hi, %s", registerDoctor.FullName),
	})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while sending code to doctor's email") {
		return
	}

	c.JSON(http.StatusOK, models.RegisterRespModel{
		Message: message,
	})
}

// Verify doctor
// @Router /v1/doctor/verify/{email}/{code} [get]
// @Summary verify doctor
// @Tags Doctor
// @Description Verify a doctor with code sent to their email
// @Accept json
// @Product json
// @Param email path string true "email"
// @Param code path string true "code"
// @Success 201 {object} models.Status
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) VerifyDoctor(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	doctorEmail := c.Param("email")
	doctorCode := c.Param("code")

	registeredDoctor, err := redis.Bytes(h.inMemoryStorage.Get(doctorEmail))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardErrorModel{
			Status:  ErrorCodeNotFound,
			Message: "Code is expired, try again",
		})
		h.log.Error("Code is expired, TTL is over.")
		return
	}

	var doctor models.RegisterDoctor
	if err := json.Unmarshal(registeredDoctor, &doctor); err != nil {
		if handleInternalServerErrorWithMessage(c, h.log, err, "cannot unmarshal doctor from redis") {
			return
		}
	}

	if doctor.Code != doctorCode {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("code is incorrect, verification is failed"), ErrorCodeInvalidCode) {
			return
		}
	}

	doctor.Password, err = etc.GenerateHashPassword(doctor.Password)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while hashing the password") {
		return
	}

	h.jwtHandler = tokens.JWTHandler{
		Sub:       doctor.ID,
		Role:      "doctor",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		TimeOut:   h.cfg.AccessTokenTimeOut,
	}

	_, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if handleInternalServerErrorWithMessage(c, h.log, err, "error generating access and refresh tokens") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err = h.serviceManager.HealthCareService().CreateDoctor(ctx, &pb.Doctor{
		Id:            doctor.ID,
		FullName:      doctor.FullName,
		Password:      doctor.Password,
		BirthDate:     doctor.BirthDate,
		Gender:        doctor.Gender,
		PhoneNumber:   doctor.PhoneNumber,
		Email:         doctor.Email,
		Address:       doctor.Address,
		Salary:        float32(doctor.Salary),
		Biography:     doctor.Biography,
		StartWorkYear: doctor.StartWorkYear,
		EndWorkYear:   doctor.EndWorkYear,
		WorkYears:     doctor.WorkYears,
		DepartmentId:  doctor.DepartmentId,
		RefreshToken:  refresh,
		IsVerified:    false,
	})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while creating doctor(is_verified = false yet)") {
		return
	}
	c.JSON(http.StatusOK, models.Status{
		Message: "registration completed successfully, wait for the admin's verification",
	})
}

// Login doctor
// @Summary login doctor
// @Tags Doctor
// @Description Login
// @Accept json
// @Produce json
// @Param User body models.LoginReqModel true "Login"
// @Success 201 {object} models.LoginRespDoctor
// @Failure 400 string Error models.ResponseError
// @Failure 400 string Error models.ResponseError
// @Router /v1/doctor/login [post]
func (h *handlerV1) LoginDoctor(c *gin.Context) {
	var (
		jspMarshal protojson.MarshalOptions
		body       models.LoginReqModel
	)

	jspMarshal.UseProtoNames = true
	err := c.ShouldBind(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, "error while marshaling doctor request body") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	doctor, err := h.serviceManager.HealthCareService().Exists(ctx, &pb.Email{Email: body.Email})

	if handleInternalServerErrorWithMessage(c, h.log, err, "error while checking if doctor exists") {
		return
	}

	if !etc.CompareHashPassword(doctor.Password, body.Password) {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("wrong password"), ErrorInvalidCredentials) {
			return
		}
	}

	h.jwtHandler = tokens.JWTHandler{
		Sub:       doctor.Password,
		Role:      "doctor",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		TimeOut:   h.cfg.AccessTokenTimeOut,
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while generating access and refresh token") {
		return
	}

	loginResp := models.LoginRespDoctor{
		ID:            doctor.Id,
		FullName:      doctor.FullName,
		Password:      doctor.Password,
		BirthDate:     doctor.BirthDate,
		Gender:        doctor.Gender,
		PhoneNumber:   doctor.PhoneNumber,
		Email:         doctor.Email,
		Address:       doctor.Address,
		Salary:        float64(doctor.Salary),
		Biography:     doctor.Biography,
		StartWorkYear: doctor.StartWorkYear,
		EndWorkYear:   doctor.EndWorkYear,
		WorkYears:     doctor.WorkYears,
		DepartmentId:  doctor.DepartmentId,
		AccessToken:   access,
		IsVerified:    doctor.IsVerified,
		CreatedAt:     doctor.CreatedAt,
		UpdatedAt:     doctor.UpdatedAt,
	}

	c.JSON(http.StatusOK, loginResp)
}

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
		IsVerified:    true,
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
		IsVerified:    respDoctor.IsVerified,
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
		IsVerified:    respDoctor.IsVerified,
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
		IsVerified:    respDoctor.IsVerified,
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

// List doctors by department id
// @Router /v1/doctors/{page}/{limit}/{department_id} [get]
// @Security BearerAuth
// @Summary get doctors' list by department id
// @Tags Doctor
// @Description get all doctors by department id
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Param department_id path int64 true "department_id"
// @Success 201 {object} models.ListDoctors
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListDoctorsByDepartmentId(c *gin.Context) {
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

	departmentId := c.Query("department_id")
	idToInt, err := strconv.Atoi(departmentId)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.HealthCareService().GetAllDoctorsByDepartmentId(ctx, &pb.GetRequest{Page: int64(page), Limit: int64(limit), Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}
