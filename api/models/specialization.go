package models

type SpecializationReq struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentId int64  `json:"department_id"`
}

type SpecializationModel struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentId int64  `json:"department_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ListSpecializations struct {
	Count           int64                  `json:"count"`
	Specializations []*SpecializationModel `json:"specializations"`
}

type ListSpecializationsWithPrices struct {
	Count           int64                     `json:"count"`
	Specializations []*SpecializaionWithPrice `json:"specializations"`
}

type SpecializaionWithPrice struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	DepartmentId int64   `json:"department_id"`
	OnlinePrice  float32 `json:"online_price"`
	OfflinePrice float32 `json:"offline_price"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
