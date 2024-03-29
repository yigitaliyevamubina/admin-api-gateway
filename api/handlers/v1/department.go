package v1

import (
	"context"
	"fmt"
	"myproject/admin-api-gateway/api/models"
	pb "myproject/admin-api-gateway/genproto/healthcare-service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

var FilePath string

// CreateDepartment
// @Router /v1/department/create [post]
// @Security BearerAuth
// @Summary create department
// @Tags Department
// @Description Create a new department with the provided details
// @Accept json
// @Produce json
// @Param DepartmentInfo body models.Department true "Create department"
// @Success 201 {object} models.DepartmentResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateDepartment(c *gin.Context) {
	var (
		body       models.Department
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	err = body.Validate()
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorValidationError) {
		return
	}

	createReq := &pb.Department{
		Name:        body.Name,
		Description: body.Description,
		ComeTime:    "2000-12-12 " + body.ComeTime + ":00",
		FinishTime:  "2000-12-12 " + body.FinishTime + ":00",
		ImageUrl:    body.ImageUrl,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respDepartment, err := h.serviceManager.HealthCareService().CreateDepartment(ctx, createReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while creating a department") {
		return
	}

	response := models.DepartmentResp{
		ID:          respDepartment.Id,
		Name:        respDepartment.Name,
		Description: respDepartment.Description,
		ComeTime:    respDepartment.ComeTime,
		FinishTime:  respDepartment.FinishTime,
		ImageUrl:    respDepartment.ImageUrl,
		CreatedAt:   respDepartment.CreatedAt,
		UpdatedAt:   respDepartment.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Get Department By Id
// @Router /v1/department/{id} [get]
// @Security BearerAuth
// @Summary get department by id
// @Tags Department
// @Description Get department
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.DepartmentResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetDepartmentById(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	respDepartment, err := h.serviceManager.HealthCareService().GetDepartmentById(ctx, &pb.GetReqInt{Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while getting department by id") {
		return
	}

	response := models.DepartmentResp{
		ID:          respDepartment.Id,
		Name:        respDepartment.Name,
		Description: respDepartment.Description,
		ComeTime:    respDepartment.ComeTime,
		FinishTime:  respDepartment.FinishTime,
		ImageUrl:    respDepartment.ImageUrl,
		CreatedAt:   respDepartment.CreatedAt,
		UpdatedAt:   respDepartment.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Update Department
// @Router /v1/department/update/{id} [put]
// @Security BearerAuth
// @Summary update department
// @Tags Department
// @Description Update department
// @Accept json
// @Produce json
// @Param id path int64 false "id"
// @Param UserInfo body models.DoctorUpdateReq true "Update Department"
// @Success 201 {object} models.DoctorResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateDepartment(c *gin.Context) {
	var (
		body        models.Department
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	updateReq := &pb.Department{
		Id:          body.ID,
		Name:        body.Name,
		Description: body.Description,
		ComeTime:    body.ComeTime,
		FinishTime:  body.FinishTime,
		ImageUrl:    body.ImageUrl,
	}

	if body.ID == 0 {
		id := c.Param("id")
		if id == "" {
			if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("id is required"), ErrorCodeInvalidParams) {
				return
			}
		}
		idToInt, err := strconv.Atoi(id)
		if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
			return
		}
		body.ID = int32(idToInt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respDepartment, err := h.serviceManager.HealthCareService().UpdateDepartment(ctx, updateReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	response := models.DepartmentResp{
		ID:          respDepartment.Id,
		Name:        respDepartment.Name,
		Description: respDepartment.Description,
		ComeTime:    respDepartment.ComeTime,
		FinishTime:  respDepartment.FinishTime,
		ImageUrl:    respDepartment.ImageUrl,
		CreatedAt:   respDepartment.CreatedAt,
		UpdatedAt:   respDepartment.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Delete Department
// @Router /v1/department/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete doctor
// @Tags Department
// @Description Delete department
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteDepartment(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err := h.serviceManager.HealthCareService().DeleteDoctor(ctx, &pb.GetReqStr{Id: id})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while deleting department") {
		return
	}

	c.JSON(http.StatusOK, models.Status{Message: "department was successfully deleted"})
}

// List departments
// @Router /v1/departments/{page}/{limit} [get]
// @Security BearerAuth
// @Summary get departments' list
// @Tags Department
// @Description get all departments
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Success 201 {object} models.ListDepartments
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListDepartments(c *gin.Context) {
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

	response, err := h.serviceManager.HealthCareService().GetAllDepartments(ctx, &pb.GetAll{Page: int64(page), Limit: int64(limit)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// Upload files
// @Router /v1/department/upload [post]
// @Summary upload department' images
// @Tags Department
// @Description upload department's images
// @Accept image/png
// @Produce json
// @Param file formData file true "file"
// @Success 201 {object} models.URL
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UploadDepartmentFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	defer file.Close()

	id := uuid.New()
	fileName := id.String() + filepath.Ext(header.Filename)
	dst, _ := os.Getwd()
	uploadPath := filepath.Join(dst, "media", "departments")

	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	FilePath = filepath.Join(uploadPath, fileName)
	err = c.SaveUploadedFile(header, FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	baseURL := "http://" + c.Request.Host
	fileURL := baseURL + "/media/departments/" + fileName

	c.JSON(http.StatusCreated, models.URL{
		URL: fileURL,
	})
}
