package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	ID       uint
	Email    string
	Password string
	Role     string
	//his refresh token will be save in httpOnly cookie
	RefreshToken string
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() *GormUserRepository {
	db, err := Connect()
	if err != nil {
		panic("lol err database connect")
	}
	return &GormUserRepository{db: db}
}

// Get from data of Account and delete by email of account
func (r *GormUserRepository) DeleteUser(id uint) {
	r.db.Where("id = ? ", id).Delete(&Account{})
	//return r.db.Create(...).Error for return err if not nil
}

func (r *GormUserRepository) AddNewUser(ac *Account) {
	r.db.Create(&Account{Email: ac.Email, Password: ac.Password, Role: ac.Role, RefreshToken: ""}) //add a new record
	//return r.db.Create(...).Error for return err if not nil
}

// check account for registration
func (r *GormUserRepository) CheckAvailibleEmail(ac *Account) bool {
	var user Account
	r.db.Where("email = ?", ac.Email).First(&user) //return nil если запись не нашлась
	if user.Email != "" {
		fmt.Println("запись найдена")
		return false
	}
	fmt.Println("запись найдeна")
	return true
}

func (r *GormUserRepository) CheckPasswordUser(ac *Account) bool {
	var user Account
	r.db.Where("email = ?", ac.Email).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ac.Password))
	return err == nil //return if equal otherwise return false
}

func (r *GormUserRepository) GetIDUser(ac *Account) uint {
	var user Account
	r.db.Where("email = ?", ac.Email).First(&user)
	return user.ID
}

func (r *GormUserRepository) SetRefreshToken(idUser, refreshToken string) {
	r.db.Model(&Account{}).Where("id = ?", idUser).UpdateColumn("refresh_token", refreshToken)
}

func (r *GormUserRepository) GetRefreshTokenUser(idUser string) string {
	var user Account
	r.db.Where("id = ?", idUser).First(&user)
	return user.RefreshToken
}
