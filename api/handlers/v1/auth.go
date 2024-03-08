package v1

import (
	"context"
	"fmt"
	"myproject/admin-api-gateway/api/handlers/tokens"
	"myproject/admin-api-gateway/api/models"
	"myproject/admin-api-gateway/pkg/etc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// Create Admin
// @Router /v1/auth/create [post]
// @Security BearerAuth
// @Summary create admin
// @Tags Auth
// @Description Create a new admin if you are a superadmin
// @Accept json
// @Product json
// @Param username query string true "username"
// @Param password query string true "password"
// @Param admin body models.AdminReq true "admin"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) CreateAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AdminReq
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")

	if superAdminUsername == "admin" && superAdminPassword == "admin" {
		err := c.BindJSON(&body)
		if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
			return
		}

		body.Id = uuid.NewString()

		body.Password, err = etc.GenerateHashPassword(body.Password)
		if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
			return
		}

		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Id,
			Role:      "admin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			TimeOut:   h.cfg.AccessTokenTimeout,
		}

		_, refresh, err := h.jwtHandler.GenerateAuthJWT()
		if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
			return
		}

		adminResp := models.AdminResp{
			Id:           body.Id,
			FullName:     body.FullName,
			Role:         body.Role,
			Age:          body.Age,
			UserName:     body.UserName,
			Email:        body.Email,
			Password:     body.Password,
			RefreshToken: refresh,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
		defer cancel()

		err = h.postgres.Create(ctx, &adminResp)
		if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
			return
		}

		c.JSON(http.StatusCreated, models.SuperAdminMessage{
			Message: "admin successfully created",
		})
	} else {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("you cannot create admin, provide username and password"), ErrorCodeInvalidJSON) {
			return
		}
	}
}

// Delete Admin
// @Router /v1/auth/delete [delete]
// @Security BearerAuth
// @Summary delete admin
// @Tags Auth
// @Description delete admin if you are a superadmin
// @Accept json
// @Product json
// @Param username query string false "username"
// @Param password query string false "password"
// @Param admin body models.DeleteAdmin true "admin"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) DeleteAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.DeleteAdmin
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")

	if superAdminUsername != "admin" && superAdminPassword != "admin" {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("you cannot create admin, provide username and password"), ErrorCodeUnauthorized) {
			return
		}
	}

	err := c.BindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	_, password, status, err := h.postgres.Check(ctx, body.Username)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	if !status {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("this admin does not exist"), ErrorCodeNotFound) {
			return
		}
	}

	if !etc.CompareHashPassword(password, body.Password) {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("incorrect password"), ErrorInvalidCredentials) {
			return
		}
	}

	resp := h.postgres.Delete(ctx, body.Username, body.Password)
	if resp != nil {
		if resp.Error() == "no rows were deleted" {
			if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("this admin does not exist"), ErrorCodeNotFound) {
				return
			}
		}
		if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
			return
		}
	}

	c.JSON(http.StatusOK, models.SuperAdminMessage{
		Message: "admin is successfully deleted",
	})
}

// Login Admin
// @Router /v1/auth/login [post]
// @Summary login
// @Tags Auth
// @Description login as admin
// @Accept json
// @Product json
// @Param admin body models.AdminLoginReq true "Login"
// @Success 201 {object} models.AdminLoginResp
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) LoginAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AdminLoginReq
	)
	jspbMarshal.UseProtoNames = true

	err := c.BindJSON(&body)

	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	role, password, status, err := h.postgres.Check(ctx, body.Username)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	if !status {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("this admin does not exist"), ErrorCodeNotFound) {
			return
		}
	}

	if !etc.CompareHashPassword(password, body.Password) {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("incorrect password"), ErrorInvalidCredentials) {
			return
		}
	}

	if role == "admin" {
		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Username,
			Role:      "admin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			TimeOut:   h.cfg.AccessTokenTimeout,
		}
	} else if role == "superadmin" {
		fmt.Println(role)
		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Username,
			Role:      "superadmin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			TimeOut:   h.cfg.AccessTokenTimeout,
		}
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	response := models.AdminLoginResp{
		Success:     true,
		AccessToken: access,
	}

	c.JSON(http.StatusOK, response)
}

// List admins
// @Router /v1/auth/admins/{page}/{limit} [get]
// @Security BearerAuth
// @Summary list admins
// @Tags Auth
// @Description list all admins
// @Accept json
// @Product json
// @Param page path string false "page"
// @Param limit path string false "limit"
// @Success 201 {object} models.ListAdminsResp
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) ListAdmins(c *gin.Context) {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	response, err := h.postgres.ListAdmins(ctx, models.ListAdminReq{Page: int32(page), Limit: int32(limit)})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get Admin
// @Router /v1/auth/get/{id} [get]
// @Security BearerAuth
// @Summary get admin
// @Tags Auth
// @Description get admin
// @Accept json
// @Product json
// @Param id path string true "id"
// @Success 201 {object} models.AdminReq
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) GetAdmin(c *gin.Context) {
	var jspMarshal protojson.MarshalOptions
	jspMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((h.cfg.CtxTimeout)))
	defer cancel()

	respAdmin, err := h.postgres.GetAdmin(ctx, models.GetAdminReq{Id: id})
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, respAdmin)
}

// Update Admin
// @Router /v1/auth/update [put]
// @Security BearerAuth
// @Summary update admin
// @Tags Auth
// @Description update admin
// @Accept json
// @Product json
// @Param admin body models.AdminUpdateReq true "admin"
// @Success 201 {object} models.AdminReq
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
func (h *handlerV1) Update(c *gin.Context) {
	var (
		body       models.AdminUpdateReq
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if handleBadRequestErrWithMessage(c, h.log, err, ErrorCodeInvalidJSON) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, _, status, err := h.postgres.Check(ctx, body.UserName)
	if handleInternalServerErrorWithMessage(c, h.log, err, ErrorCodeInternalServerError) {
		return
	}

	if status {
		if handleBadRequestErrWithMessage(c, h.log, fmt.Errorf("this is username is used by another admin, try a new username"), ErrorBadRequest) {
			return
		}
	}

	response, err := h.postgres.Update(ctx, &body)
	if handleInternalServerErrorWithMessage(c, h.log, err, "failed to update admin") {
		return
	}

	c.JSON(http.StatusOK, response)
}
