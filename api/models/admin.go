package models

type AdminReq struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Age      int64  `json:"age"`
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AdminUpdateReq struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Age      int64  `json:"age"`
	UserName string `json:"username"`
}

type AdminResp struct {
	Id           string `json:"id"`
	FullName     string `json:"full_name"`
	Age          int64  `json:"age"`
	Email        string `json:"email"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
}

type AdminLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResp struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
}

type SuperAdminMessage struct {
	Message string `json:"message"`
}

type DeleteAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RbacAllRolesResp struct {
	Roles []string `json:"roles"`
}

type Policy struct {
	Role     string `json:"role"`
	EndPoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type ListRolePolicyResp struct {
	Policies []*Policy `json:"policies"`
}

type AddPolicyRequest struct {
	Policy Policy `json:"policy"`
}

type ListAdminReq struct {
	Page  int32
	Limit int32
}

type ListAdminsResp struct {
	Count  int64       `json:"count"`
	Admins []*AdminReq `json:"admins"`
}

type GetAdminReq struct {
	Id string
}
