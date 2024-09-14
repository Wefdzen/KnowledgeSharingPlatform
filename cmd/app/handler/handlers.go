package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"wefdzen/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput database.Account
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		tmp := database.Account{
			Email:    jsonInput.Email,
			Password: jsonInput.Password,
		}
		//Check password account
		userRepo := database.NewGormUserRepository()
		result := database.CheckPasssword(userRepo, &tmp)
		if !result { //password not equal
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "incorrect password",
			})
			return
		}
		//get id for jwt
		idUser := database.GetID(userRepo, &tmp)
		//id будет его айди из бд думай
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":        idUser,                                  //id from bd
			"liveToken": time.Now().Add(time.Minute * 15).Unix(), // 30 min
			"role":      jsonInput.Role,                          //admin, user
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("secret_key")))
		if err != nil {
			panic("lol access token")
		}

		token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"liveToken": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 day
		})
		tokenString2, err := token2.SignedString([]byte(os.Getenv("secret_key")))
		if err != nil {
			panic("lol refresh token")
		}
		//set refreshToken in db
		userId := fmt.Sprintf("%v", idUser) // Преобразование id в строку
		database.SetRefToken(userRepo, userId, tokenString2)
		//setcookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", tokenString, 3600*24*30, "", "", false, true) // хранится тоже 30 дней в браузере ->
		// -> потому что если он продет в middleware он не перейдет в проверку refresh
		c.SetCookie("refreshToken", tokenString2, 3600*24*30, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{})
	}
}

func Registration() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get data from json
		var jsonInput database.Account
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		//hashPassword
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(jsonInput.Password), 12)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		jsonInput.Password = string(hashPassword)

		tmp := database.Account{
			Email:        jsonInput.Email,
			Password:     jsonInput.Password,
			Role:         "user",
			RefreshToken: "",
		}
		//connect
		userRepo := database.NewGormUserRepository()
		//проверка на наличие такой почты
		result := database.EmailAvailible(userRepo, &tmp)
		if !result {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email is busy",
			})
			return
		} // success

		//registration in db
		database.RegisterUser(userRepo, &tmp)

		c.JSON(http.StatusOK, gin.H{
			"status": "create succes",
		})
	}
}

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepo := database.NewGormUserRepository()
		database.RemoveUser(userRepo, 3) //get id from jwt
	}
}

func MainPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access": "success",
		})
	}
}
