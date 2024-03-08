package middleware

import (
	"myproject/admin-api-gateway/api/handlers/tokens"
	v1 "myproject/admin-api-gateway/api/handlers/v1"
	"myproject/admin-api-gateway/api/models"
	"myproject/admin-api-gateway/config"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type CasbinHandler struct {
	cfg      config.Config
	enforcer *casbin.Enforcer
}

func Auth(casbin *casbin.Enforcer, cfg config.Config) gin.HandlerFunc {
	casbinHandler := &CasbinHandler{
		cfg:      cfg,
		enforcer: casbin,
	}

	return func(ctx *gin.Context) {
		allowed, err := casbinHandler.CheckPermission(ctx.Request)
		if err != nil {
			v, _ := err.(jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				casbinHandler.RequireRefresh(ctx)
				return
			} else {
				casbinHandler.RequirePermission(ctx)
			}
		} else if !allowed {
			casbinHandler.RequirePermission(ctx)
			return
		}
	}
}

func (c *CasbinHandler) GetRole(ctx *http.Request) (string, int) {
	token := ctx.Header.Get("Authorization")
	if token == "" {
		return "unauthorized", http.StatusOK
	}

	var cutToken string
	if strings.Contains(token, "Bearer") {
		cutToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		cutToken = token
	}

	claims, err := tokens.ExtractClaims(cutToken, []byte(c.cfg.SignInKey))
	if err != nil {
		return "unauthorized, token is invalid", http.StatusBadRequest
	}

	return cast.ToString(claims["role"]), http.StatusOK
}

func (c *CasbinHandler) CheckPermission(ctx *http.Request) (bool, error) {
	role, status := c.GetRole(ctx)
	if status != http.StatusOK {
		return false, nil
	}

	object := ctx.URL.Path
	action := ctx.Method

	response, err := c.enforcer.Enforce(role, object, action)
	if err != nil {
		return false, nil
	}

	return response, nil
}

func (c *CasbinHandler) RequirePermission(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, models.StandardErrorModel{
		Status:  v1.StatusMethodNotAllowed,
		Message: "This method is not allowed to you",
	})
}

func (c *CasbinHandler) RequireRefresh(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.StandardErrorModel{
		Status:  v1.ErrorCodeUnauthorized,
		Message: "Access token is expired, refresh it.",
	})
}
