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
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token payload",
			})
		}

		userID := claims["user_id"].(string)
		c.Locals("user_id", userID)

		// 🔧 Автосоздание профиля, если нет
		var profile model.UserProfile
		if err := db.First(&profile, "user_id = ?", userID).Error; err == gorm.ErrRecordNotFound {
			db.Create(&model.UserProfile{
				UserID:    userID,
				Nickname:  "Okkura chan fan",
				FirstName: "",
				LastName:  "",
				AvatarURL: "",
			})
		}

		return c.Next()
	}
}
