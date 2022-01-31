package models

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	GenderID  int16  `json:"gender_id"`
	RoleID    int16  `json:"role_id"`
	Coin      int16  `json:"coin"`
	IsDeleted bool   `json:"is_deleted"`
	CreateAt  string `json:"create_at"`
	UpdateAt  string `json:"update_at"`
	DeleteAt  string `json:"delete_at"`
	Image     string `json:"image"`
}

type Update struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	GenderID int16  `json:"gender_id" validate:"required"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

type UserRequest struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	AdminID int64  `json:"admin_id"`
	IsAcc   bool   `json:"is_acc"`
	Request string `json:"request_date"`
	Review  string `json:"review_date"`
}

type Request struct {
	UserID int64 `json:"user_id"`
	Answer bool  `json:"answer"`
}
