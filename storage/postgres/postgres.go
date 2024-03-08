package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myproject/admin-api-gateway/api/models"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) Create(ctx context.Context, admin *models.AdminResp) error {
	query := `INSERT INTO admins(id, full_name, age, username, email, password, role, refresh_token) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, admin.Id,
		admin.FullName,
		admin.Age,
		admin.UserName,
		admin.Email,
		admin.Password,
		admin.Role,
		admin.RefreshToken)

	return err
}

func (r *adminRepo) Update(ctx context.Context, adminReq *models.AdminUpdateReq) (*models.AdminReq, error) {
	query := `UPDATE admins SET 
					full_name = $1, 
					age = $2, 
					username = $3 
					WHERE id = $4 
					RETURNING id, 
							  full_name, 
							  age, 
							  username, 
							  email, 
							  password, 
							  role`
	row := r.db.QueryRow(query, adminReq.FullName, adminReq.Age, adminReq.UserName, adminReq.Id)

	var admin models.AdminReq
	if err := row.Scan(&admin.Id,
		&admin.FullName,
		&admin.Age,
		&admin.UserName,
		&admin.Email,
		&admin.Password,
		&admin.Role); err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *adminRepo) Delete(ctx context.Context, userName, password string) error {
	query := `DELETE FROM admins WHERE username = $1`
	result, err := r.db.Exec(query, userName)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("error")
		return errors.New("no rows were deleted")
	}

	return nil
}

func (r *adminRepo) Check(ctx context.Context, userName string) (string, string, bool, error) {
	query := `SELECT COUNT(1), password, role
	FROM admins GROUP by username, password, role
	HAVING username = $1
	`
	var status int
	var password string
	var role string
	result := r.db.QueryRow(query, userName)
	if err := result.Scan(&status, &password, &role); err != nil {
		return "", "", false, nil
	}

	return role, password, status == 1, nil
}

func (r *adminRepo) GetAdmin(rctx context.Context, req models.GetAdminReq) (*models.AdminReq, error) {
	query := `SELECT id, 
				full_name, 
				age, 
				username, 
				email, 
				password, 
				role FROM admins WHERE id = $1`

	row := r.db.QueryRow(query, req.Id)

	var admin models.AdminReq
	if err := row.Scan(&admin.Id,
		&admin.FullName,
		&admin.Age,
		&admin.UserName,
		&admin.Email,
		&admin.Password,
		&admin.Role); err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *adminRepo) ListAdmins(ctx context.Context, req models.ListAdminReq) (*models.ListAdminsResp, error) {
	query := `SELECT id, 
				full_name, 
				age, 
				username, 
				email, 
				password, 
				role FROM admins`

	var admins models.ListAdminsResp

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var admin models.AdminReq
		if err := rows.Scan(&admin.Id,
			&admin.FullName,
			&admin.Age,
			&admin.UserName,
			&admin.Email,
			&admin.Password,
			&admin.Role); err != nil {
			return nil, err
		}
		admins.Admins = append(admins.Admins, &admin)
		admins.Count++
	}

	return &admins, nil
}
