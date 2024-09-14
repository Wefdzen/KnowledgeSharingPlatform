package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"wefdzen/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем access token из cookie
		accessToken, err := c.Cookie("accessToken")
		if err != nil { // если не удалось получить токен
			c.AbortWithStatus(http.StatusUnauthorized)
			return // необходимо прекратить выполнение
		}

		// Парсим токен без проверки exp
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Возвращаем секретный ключ
			return []byte(os.Getenv("secret_key")), nil
		})

		if err != nil {
			fmt.Println("Ошибка парсинга токена:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Получаем claims из токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Проверяем вручную exp = liveToken поле Parse сам проверяет exp
			if float64(time.Now().Unix()) >= claims["liveToken"].(float64) {
				// Если токен истек, проверяем refreshToken
				userId := fmt.Sprintf("%v", claims["id"]) // Преобразование id в строку
				if !CheckRefreshToken(c, userId) {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				} else {
					fmt.Println("Ты зашел с помощью refresh token")
					//Создать новый access token and refresh и отправить его в cookie

					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"id":        claims["id"],                            //id from bd
						"liveToken": time.Now().Add(time.Minute * 15).Unix(), // 15 min
						"role":      claims["role"],                          //admin, user
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
					//connect
					userRepo := database.NewGormUserRepository()
					//update refreshToken in db
					database.SetRefToken(userRepo, userId, tokenString2)
					//setcookie
					c.SetSameSite(http.SameSiteLaxMode)
					c.SetCookie("accessToken", tokenString, 3600*24*30, "", "", false, true) // хранится тоже 30 дней в браузере ->
					// -> потому что если он продет в middleware он не перейдет в проверку refresh
					c.SetCookie("refreshToken", tokenString2, 3600*24*30, "", "", false, true)
				}
			} else {
				fmt.Println("Ты зашел с помощью access token")
			}
		} else {
			// Если claims не валидны
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next() // продолжаем обработку запроса
	}
}

// я понимаю что тут x2 проверка
func CheckRefreshToken(c *gin.Context, idUser string) bool {
	// Получаем refresh token из cookie
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil { // если не удалось получить токен
		c.AbortWithStatus(http.StatusUnauthorized)
		return false
	}

	// Проверяем refresh token в базе данных
	userRepo := database.NewGormUserRepository()
	TokenForCheck := database.GetRefToken(userRepo, idUser)

	if TokenForCheck != refreshToken {
		return false
	}

	// Парсим refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret_key")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return false
	}

	// Если refresh token истек, возвращаем false
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) >= claims["liveToken"].(float64) {
			return false
		}
	} else {
		return false
	}

	return true
}
