package handler

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"donate-service/internal/model"
)

func GetMyDonates(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		var donations []model.Donation
		if err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&donations).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Не удалось получить донаты"})
		}

		return c.JSON(donations)
	}
}

func GetMySummary(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		var count int64
		var total sql.NullInt64

		// Считаем количество донатов
		db.Model(&model.Donation{}).
			Where("user_id = ?", userID).
			Count(&count)

		// Суммируем сумму донатов
		db.Model(&model.Donation{}).
			Where("user_id = ?", userID).
			Select("SUM(amount)").
			Scan(&total)

		sum := int64(0)
		if total.Valid {
			sum = total.Int64
		}

		// Определяем подписку
		subscription := "—"
		if sum >= 999 {
			subscription = "Ультимативный"
		} else if sum >= 599 {
			subscription = "Премиум"
		} else if sum >= 299 {
			subscription = "Базовый"
		}

		return c.JSON(fiber.Map{
			"donateCount":  count,
			"totalAmount":  sum,	
			"subscription": subscription,
		})
	}
}


func GetProfile(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		// Получаем профиль
		var profile model.UserProfile
		if err := db.First(&profile, "user_id = ?", userID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Profile not found",
			})
		}

		// Получаем email и роль из основной таблицы users
		var email, role string
		err := db.Raw("SELECT email, role FROM users WHERE id = ?", userID).Row().Scan(&email, &role)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to load user info",
			})
		}

return c.JSON(fiber.Map{
	"user_id":    userID,
	"nickname":   profile.Nickname,
	"first_name": profile.FirstName,
	"last_name":  profile.LastName,
	"avatar_url": profile.AvatarURL,
})

	}
}



func UpdateProfile(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		var input struct {
			Nickname  string `json:"nickname"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			AvatarURL string `json:"avatar_url"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid input",
			})
		}

		var profile model.UserProfile
		if err := db.First(&profile, "user_id = ?", userID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "profile not found",
			})
		}

		profile.Nickname = input.Nickname
		profile.FirstName = input.FirstName
		profile.LastName = input.LastName
		profile.AvatarURL = input.AvatarURL

		if err := db.Save(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update",
			})
		}

		return c.JSON(profile)
	}
}
