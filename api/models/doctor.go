package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

type DoctorReq struct {
	ID            string  `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	BirthDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
	SpecIds       []int64 `json:"spec_ids"`
}

type RegisterDoctor struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
	SpecIds       []int64 `json:"spec_ids"`
	IsVerified    bool    `json:"is_verified"`
	Code          string  `json:"code"`
}

func (d *DoctorReq) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.FirstName, validation.Required, validation.Length(3, 50).Error("the length should be 3-50 characters"), validation.Match(regexp.MustCompile(`^[A-Z][a-z]*$`)).Error("should start with a capital letter and should only contain letters")),
		validation.Field(&d.LastName, validation.Required, validation.Length(3, 50).Error("the length should be 3-50 characters"), validation.Match(regexp.MustCompile(`^[A-Z][a-z]*$`)).Error("should start with a capital letter and should only contain letters")),
		validation.Field(&d.BirthDate, validation.Required, validation.Match(regexp.MustCompile(`^\\d{4}-\\d{2}-\\d{2}$`)).Error("should be in the format 'yyyy-mm-dd'")),
		validation.Field(&d.Gender, validation.Required, validation.Length(4, 6).Error("the length should be 4-6 characters"), validation.In("male", "female").Error("should either be male or female only")),
		validation.Field(&d.Email, validation.Required, is.Email),
		validation.Field(&d.StartWorkYear, validation.Required, validation.Match(regexp.MustCompile(`^\\d{4}-\\d{2}-\\d{2}$`)).Error("should be in the format 'yyyy-mm-dd'")),
		validation.Field(&d.EndWorkYear, validation.Required, validation.Match(regexp.MustCompile(`^\\d{4}-\\d{2}-\\d{2}$`)).Error("should be in the format 'yyyy-mm-dd'")),
		validation.Field(&d.Password,
			validation.Required,
			validation.Length(5, 30).Error("the length should ve 5-30 characters"),
			validation.Match(regexp.MustCompile(`\\d`)).Error(`should contain at least one digit`),
			validation.Match(regexp.MustCompile(`^[a-zA-Z\\d]+$`)).Error(`should only contain letters (either lowercase or uppercase) and digits`),
		),
		validation.Field(&d.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(`^\d{9}$`))),
	)
}

type DoctorModel struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
	AccessToken   string  `json:"access_token"`
	SpecIds       []int64 `json:"spec_ids"`
	IsVerified    bool    `json:"is_verified"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type DoctorResp struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
	SpecIds       []int64 `json:"spec_ids"`
	IsVerified    bool    `json:"is_verified"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type LoginRespDoctor struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
	SpecIds       []int64 `json:"spec_ids"`
	IsVerified    bool    `json:"is_verified"`
	AccessToken   string  `json:"access_token"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ListReq struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type ListDoctors struct {
	Count   int64         `json:"count"`
	Doctors []*DoctorResp `json:"doctors"`
}

type DoctorUpdateReq struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	Address       string  `json:"address"`
	Salary        float64 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     int64   `json:"work_years"`
	DepartmentId  int64   `json:"department_id"`
}
