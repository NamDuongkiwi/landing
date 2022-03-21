package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id          int64  `json:"id"`
	Email       int64  `json:"email"`
	Name        string `json:"name"`
	Password    []byte `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"'updated_at'"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
