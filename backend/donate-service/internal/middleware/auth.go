package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"donate-service/internal/model"
)

func RequireAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем заголовок авторизации
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		// Извлекаем токен
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		// Проверяем токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Извлекаем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token payload",
			})
		}

		userID := claims["user_id"].(string)
		c.Locals("user_id", userID)

		// Автосоздание профиля, если его нет
		var profile model.UserProfile
		if err := db.First(&profile, "user_id = ?", userID).Error; err == gorm.ErrRecordNotFound {
			// Получаем email из токена, если есть
			email, _ := claims["email"].(string)
			
			newProfile := model.UserProfile{
				UserID:    userID,
				Email:     email,
				Nickname:  "Okkura chan fan",
				Role:      "user",
				FirstName: "",
				LastName:  "",
				AvatarURL: "",
			}

			if err := db.Create(&newProfile).Error; err != nil {
				// Логируем ошибку, но позволяем запросу продолжиться
				c.Locals("profile_create_error", err.Error())
			}
		}

		return c.Next()
	}
}
