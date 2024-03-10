package v1

import (
	"context"
	"fmt"
	"myproject/admin-api-gateway/api/models"
	pb "myproject/admin-api-gateway/genproto/healthcare-service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateSpecialization
// @Router /v1/specialization/create [post]
// @Security BearerAuth
// @Summary create specialization
// @Tags Specialization
// @Description Create a new specialization with the provided details
// @Accept json
// @Produce json
// @Param SpecializaionInfo body models.SpecializationReq true "Create specialization"
// @Success 201 {object} models.SpecializationModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateSpecializaion(c *gin.Context) {
	var (
		body       models.SpecializationReq
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	createReq := &pb.Specializations{
		Name:         body.Name,
		Description:  body.Description,
		DepartmentId: body.DepartmentId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respSpec, err := h.serviceManager.HealthCareService().CreateSpecialization(ctx, createReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while creating a specialization") {
		return
	}

	response := models.SpecializationModel{
		ID:           respSpec.Id,
		Name:         respSpec.Name,
		Description:  respSpec.Description,
		DepartmentId: respSpec.DepartmentId,
		CreatedAt:    respSpec.CreatedAt,
		UpdatedAt:    respSpec.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Get Specialization By Id
// @Router /v1/specialization/{id} [get]
// @Security BearerAuth
// @Summary get specialization by id
// @Tags Specialization
// @Description Get specialization
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.SpecializationModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetSpecializationById(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	respSpec, err := h.serviceManager.HealthCareService().GetSpecializationById(ctx, &pb.GetReqInt{Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while getting specialization by id") {
		return
	}

	response := models.SpecializationModel{
		ID:           respSpec.Id,
		Name:         respSpec.Name,
		Description:  respSpec.Description,
		DepartmentId: respSpec.DepartmentId,
		CreatedAt:    respSpec.CreatedAt,
		UpdatedAt:    respSpec.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Update Specialization
// @Router /v1/specialization/update/{id} [put]
// @Security BearerAuth
// @Summary update specialization
// @Tags Specialization
// @Description Update specialization
// @Accept json
// @Produce json
// @Param id path int64 false "id"
// @Param UserInfo body models.SpecializationReq true "Update specialization"
// @Success 201 {object} models.SpecializationModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateSpecialization(c *gin.Context) {
	var (
		body        models.SpecializationReq
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	updateReq := &pb.Specializations{
		Id:           body.ID,
		Name:         body.Name,
		Description:  body.Description,
		DepartmentId: body.DepartmentId,
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
		body.ID = int64(idToInt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respSpec, err := h.serviceManager.HealthCareService().UpdateSpecialization(ctx, updateReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	response := models.SpecializationModel{
		ID:           respSpec.Id,
		Name:         respSpec.Name,
		Description:  respSpec.Description,
		DepartmentId: respSpec.DepartmentId,
		CreatedAt:    respSpec.CreatedAt,
		UpdatedAt:    respSpec.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Delete Specialization
// @Router /v1/specialization/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete specialization
// @Tags Specialization
// @Description Delete specialization
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteSpecialization(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err = h.serviceManager.HealthCareService().DeleteSpecialization(ctx, &pb.GetReqInt{Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while deleting specialization") {
		return
	}

	c.JSON(http.StatusOK, models.Status{Message: "specialization was successfully deleted"})
}

// List specializations
// @Router /v1/specializations/{page}/{limit} [get]
// @Security BearerAuth
// @Summary get specializations' list
// @Tags Specialization
// @Description get all specializations
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Success 201 {object} models.ListSpecializations
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListSpecializations(c *gin.Context) {
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

	response, err := h.serviceManager.HealthCareService().GetAllSpecializations(ctx, &pb.GetAll{Page: int64(page), Limit: int64(limit)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// List specializations with prices by department id
// @Router /v1/specializations/{page}/{limit}/{department_id} [get]
// @Security BearerAuth
// @Summary get specializations' list by department id with prices
// @Tags Specialization
// @Description get all specializations by department id with prices
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Param department_id path int64 true "department_id"
// @Success 201 {object} models.ListSpecializationsWithPrices
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListSpecializationsByDepartmentId(c *gin.Context) {
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

	response, err := h.serviceManager.HealthCareService().GetAllSpecByDepartmentIdWithPrices(ctx, &pb.GetRequest{Page: int64(page), Limit: int64(limit), Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}
