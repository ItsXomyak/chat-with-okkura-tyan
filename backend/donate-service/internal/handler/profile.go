package handler

import (
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
		var total float64

		// Считаем количество донатов
		db.Model(&model.Donation{}).
			Where("user_id = ?", userID).
			Count(&count)

		// Суммируем сумму донатов
		db.Model(&model.Donation{}).
			Where("user_id = ?", userID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&total)

		// Определяем подписку
		subscription := "—"
		if total >= 999 {
			subscription = "Ультимативный"
		} else if total >= 599 {
			subscription = "Премиум"
		} else if total >= 299 {
			subscription = "Базовый"
		}

		return c.JSON(fiber.Map{
			"donateCount":  count,
			"totalAmount":  total,
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
			if err == gorm.ErrRecordNotFound {
				// Если профиль не найден, создаем новый
				profile = model.UserProfile{
					UserID: userID,
				}
				if err := db.Create(&profile).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to create profile",
					})
				}
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to load profile",
				})
			}
		}

		return c.JSON(fiber.Map{
			"user_id":    userID,
			"nickname":   profile.Nickname,
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
			"avatar_url": profile.AvatarURL,
			"email":      profile.Email,
			"role":       profile.Role,
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
