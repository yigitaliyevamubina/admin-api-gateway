package models

type SpecPriceReq struct {
	ID               int64   `json:"id"`
	DoctorId         string  `json:"doctor_id"`
	SpecializationId int64   `json:"specialization_id"`
	OnlinePrice      float32 `json:"online_price"`
	OfflinePrice     float32 `json:"offline_price"`
}

type SpecPriceModel struct {
	ID               int64   `json:"id"`
	DoctorId         string  `json:"doctor_id"`
	SpecializationId int64   `json:"specialization_id"`
	OnlinePrice      float32 `json:"online_price"`
	OfflinePrice     float32 `json:"offline_price"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type ListSpecPrices struct {
	Count      int64             `json:"count"`
	SpecPrices []*SpecPriceModel `json:"specialization_prices"`
}
