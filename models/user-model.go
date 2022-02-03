package models

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	GenderID uint16 `json:"gender_id"`
	RoleID   uint16 `json:"role_id"`
	Coin     uint16 `json:"coin"`
	IsActive bool   `json:"is_active"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
	DeleteAt string `json:"delete_at"`
	Image    string `json:"image"`
	// Token    string `json:"token"`
}

type Update struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	GenderID uint16 `json:"gender_id" validate:"required"`
	Image    string `json:"image"`
}

type ChangePass struct {
	OldPass    string `json:"old_pass"`
	NewPass    string `json:"new_pass"`
	RetypePass string `json:"retype_pass"`
}

type UserRequest struct {
	ID      uint64 `json:"id"`
	UserID  uint64 `json:"user_id"`
	AdminID uint64 `json:"admin_id"`
	IsAcc   bool   `json:"is_acc"`
	Request string `json:"request_date"`
	Review  string `json:"review_date"`
}

type Request struct {
	UserID uint64 `json:"user_id"`
	Answer bool   `json:"answer"`
}
