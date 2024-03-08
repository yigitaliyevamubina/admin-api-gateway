package postgresrepo

import (
	"context"
	"myproject/admin-api-gateway/api/models"
)

type AdminStorageI interface {
	Create(ctx context.Context, admin *models.AdminResp) error
	Delete(ctx context.Context, userName, password string) error
	Check(ctx context.Context, userName string) (string, string, bool, error)
	ListAdmins(ctx context.Context, req models.ListAdminReq) (*models.ListAdminsResp, error)
	GetAdmin(ctx context.Context, req models.GetAdminReq) (*models.AdminReq, error)
	Update(ctx context.Context, adminReq *models.AdminUpdateReq) (*models.AdminReq, error)
}
