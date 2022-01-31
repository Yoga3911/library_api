package models

type Login struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,max=100"`
	GenderID  int16  `json:"gender_id"`
	RoleID    int16  `json:"role_id"`
	Coin      int16  `json:"coin"`
	IsDeleted bool   `json:"is_deleted"`
	Token     string `json:"token"`
	Image     string `json:"image"`
}

type Register struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	GenderID int16  `json:"gender_id" validate:"required"`
	RoleID   int16  `json:"role_id"`
	Coin     int16  `json:"coin"`
	Token    string `json:"token"`
}

type RegisterVal struct {
	Count  int64 `json:"count"`
	EmailC int   `json:"emailC"`
	NameC  int   `json:"nameC"`
}
