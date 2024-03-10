package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

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

func (d *Department) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.ComeTime, validation.Required, validation.Match(regexp.MustCompile(`^\d{2}:\d{2}$`)).Error("should be in the format 'hh:mm'")),
		validation.Field(&d.FinishTime, validation.Required, validation.Match(regexp.MustCompile(`^\d{2}:\d{2}$`)).Error("should be in the format 'hh:mm'")),
	)
}

type ListDepartments struct {
	Count       int64             `json:"count"`
	Departments []*DepartmentResp `json:"departments"`
}
