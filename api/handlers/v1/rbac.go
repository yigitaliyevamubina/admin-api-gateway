package v1

import (
	"myproject/admin-api-gateway/api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// List all roles
// @Router /v1/rbac/roles [get]
// @Security BearerAuth
// @Summary list roles
// @Tags RBAC
// @Description List roles
// @Accept json
// @Product json
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 201 {object} models.RbacAllRolesResp
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) ListRoles(c *gin.Context) {
	var (
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true

	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")
	if superAdminPassword == "superadminpass" && superAdminUsername == "superadmin" {
		roles := h.casbin.GetAllRoles()
		c.JSON(http.StatusOK, models.RbacAllRolesResp{
			Roles: roles,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot get all roles, provide correct username and password",
		})
	}
}

// List all policies of a role
// @Router /v1/rbac/policies/{role} [get]
// @Security BearerAuth
// @Summary get all policies of a role
// @Tags RBAC
// @Description Get all policies of a role
// @Accept json
// @Product json
// @Param username query string true "username"
// @Param password query string true "password"
// @Param role path string true "role"
// @Success 201 {object} models.ListRolePolicyResp
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) ListRolePolicies(c *gin.Context) {
	var (
		jspMarshal protojson.MarshalOptions
	)

	jspMarshal.UseProtoNames = true
	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")
	if superAdminPassword == "admin" && superAdminUsername == "admin" {
		role := c.Param("role")
		var response models.ListRolePolicyResp
		for _, p := range h.casbin.GetFilteredPolicy(0, role) {
			response.Policies = append(response.Policies, &models.Policy{
				Role:     p[0],
				EndPoint: p[1],
				Method:   p[2],
			})
		}
		c.JSON(http.StatusOK, response)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot get policies of a role, provide correct username and password",
		})
	}
}

// Add policy to a role
// @Router /v1/rbac/add/policy [post]
// @Security BearerAuth
// @Summary add policy to a role
// @Tags RBAC
// @Description Add policy to a role
// @Accept json
// @Product json
// @Param username query string true "username"
// @Param password query string true "password"
// @Param policy body models.AddPolicyRequest true "policy"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) AddPolicyToRole(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AddPolicyRequest
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")
	if superAdminPassword == "admin" && superAdminUsername == "admin" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "failed, try again",
			})
			return
		}
		body.Policy.Method = strings.ToUpper(body.Policy.Method)
		status := checkMethod(body.Policy.Method)
		if !status {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   ErrorBadRequest,
				"message": "invalid method",
			})
			return
		}
		p := []string{body.Policy.Role, body.Policy.EndPoint, body.Policy.Method}
		if _, err := h.casbin.AddPolicy(p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "failed, try again",
			})
			return
		}
		h.casbin.SavePolicy()
		c.JSON(http.StatusOK, models.SuperAdminMessage{
			Message: "success",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot add policy, provide correct username and password",
		})
	}
}

// Delete policy
// @Router /v1/rbac/delete/policy [delete]
// @Security BearerAuth
// @Summary delete policy
// @Tags RBAC
// @Description Delete policy
// @Accept json
// @Product json
// @Param username query string true "username"
// @Param password query string true "password"
// @Param policy body models.AddPolicyRequest true "policy"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) DeletePolicy(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AddPolicyRequest
	)
	jspbMarshal.UseProtoNames = true
	superAdminUsername := c.Query("username")
	superAdminPassword := c.Query("password")
	if superAdminPassword == "admin" && superAdminUsername == "admin" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error,
				"message": "failed, try again",
			})
			return
		}
		p := []string{body.Policy.Role, body.Policy.EndPoint, body.Policy.Method}
		if _, err := h.casbin.RemovePolicy(p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error,
				"message": "failed, try again",
			})
			return
		}
		h.casbin.SavePolicy()
		c.JSON(http.StatusOK, models.SuperAdminMessage{
			Message: "success",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot delete policy, provide correct username and password",
		})
	}
}