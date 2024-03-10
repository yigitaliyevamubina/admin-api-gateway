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

// CreateSpecializationPrice
// @Router /v1/specprice/create [post]
// @Security BearerAuth
// @Summary create specialization price
// @Tags Specialization price
// @Description Create a new specialization price with the provided details
// @Accept json
// @Produce json
// @Param SpecPriceInfo body models.SpecPriceReq true "Create specialization price"
// @Success 201 {object} models.SpecPriceModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateSpecPrice(c *gin.Context) {
	var (
		body       models.SpecPriceReq
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	createReq := &pb.DocSpecPrices{
		DoctorId:         body.DoctorId,
		SpecializationId: body.SpecializationId,
		OnlinePrice:      body.OnlinePrice,
		OfflinePrice:     body.OfflinePrice,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respSpecPrice, err := h.serviceManager.HealthCareService().CreateDocSpecPrices(ctx, createReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while creating a specialization price") {
		return
	}

	response := models.SpecPriceModel{
		ID:               respSpecPrice.Id,
		DoctorId:         respSpecPrice.DoctorId,
		SpecializationId: respSpecPrice.SpecializationId,
		OnlinePrice:      respSpecPrice.OnlinePrice,
		OfflinePrice:     respSpecPrice.OfflinePrice,
		CreatedAt:        respSpecPrice.CreatedAt,
		UpdatedAt:        respSpecPrice.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Get Specialization Price By Id
// @Router /v1/specprice/{id} [get]
// @Security BearerAuth
// @Summary get specialization price by id
// @Tags Specialization price
// @Description Get Specialization price
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.SpecPriceModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetSpecPriceById(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	respSpecPrice, err := h.serviceManager.HealthCareService().GetSpecPriceById(ctx, &pb.GetReqInt{Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while getting specializatio price by id") {
		return
	}

	response := models.SpecPriceModel{
		ID:               respSpecPrice.Id,
		DoctorId:         respSpecPrice.DoctorId,
		SpecializationId: respSpecPrice.SpecializationId,
		OnlinePrice:      respSpecPrice.OnlinePrice,
		OfflinePrice:     respSpecPrice.OfflinePrice,
		CreatedAt:        respSpecPrice.CreatedAt,
		UpdatedAt:        respSpecPrice.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Update Specialization Price
// @Router /v1/specprice/update/{id} [put]
// @Security BearerAuth
// @Summary update specialization Price
// @Tags Specialization price
// @Description Update specialization price
// @Accept json
// @Produce json
// @Param id path int64 false "id"
// @Param UserInfo body models.SpecPriceReq true "Update specialization price"
// @Success 201 {object} models.SpecPriceModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateSpecPrice(c *gin.Context) {
	var (
		body        models.SpecPriceReq
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	updateReq := &pb.DocSpecPrices{
		Id:               body.ID,
		DoctorId:         body.DoctorId,
		SpecializationId: body.SpecializationId,
		OnlinePrice:      body.OnlinePrice,
		OfflinePrice:     body.OfflinePrice,
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

	respSpecPrice, err := h.serviceManager.HealthCareService().UpdateSpecPrice(ctx, updateReq)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	response := models.SpecPriceModel{
		ID:               respSpecPrice.Id,
		DoctorId:         respSpecPrice.DoctorId,
		SpecializationId: respSpecPrice.SpecializationId,
		OnlinePrice:      respSpecPrice.OnlinePrice,
		OfflinePrice:     respSpecPrice.OfflinePrice,
		CreatedAt:        respSpecPrice.CreatedAt,
		UpdatedAt:        respSpecPrice.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Delete Specialization Price
// @Router /v1/specprice/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete specialization price
// @Tags Specialization price
// @Description Delete specialization Price
// @Accept json
// @Produce json
// @Param id path int64 true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteSpecPrice(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidParams) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err = h.serviceManager.HealthCareService().DeleteSpecPrice(ctx, &pb.GetReqInt{Id: int64(idToInt)})
	if handleInternalServerErrorWithMessage(c, h.log, err, "error while deleting specialization price") {
		return
	}

	c.JSON(http.StatusOK, models.Status{Message: "specialization price was successfully deleted"})
}

// List specialization prices
// @Router /v1/specprices/{page}/{limit} [get]
// @Security BearerAuth
// @Summary get specialization prices' list
// @Tags Specialization price
// @Description get all specialization prices
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Success 201 {object} models.ListSpecPrices
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListSpecPrices(c *gin.Context) {
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

	response, err := h.serviceManager.HealthCareService().GetAllSpecPrice(ctx, &pb.GetAll{Page: int64(page), Limit: int64(limit)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}
