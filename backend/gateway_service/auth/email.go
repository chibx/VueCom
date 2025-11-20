package auth

import (
	"gorm.io/gorm"
)

func CheckIfEmailExist(db gorm.DB, email string) bool {
	return true
}
