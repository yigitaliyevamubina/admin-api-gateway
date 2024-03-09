package models

type Department struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ComeTime    string `json:"work_starts_at"`
	FinishTime  string `json:"work_ends_at"`
	ImageUrl    string `json:"image_url"`
}

type DepartmentResp struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ComeTime    string `json:"work_starts_at"`
	FinishTime  string `json:"work_ends_at"`
	ImageUrl    string `json:"image_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
