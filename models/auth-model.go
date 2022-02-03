package models

type Login struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type Register struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	GenderID uint16 `json:"gender_id" validate:"required"`
}

type OTP struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
