package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"donate-service/internal/middleware"
	"donate-service/internal/model"
)

/*
GET /me/profile	Получить профиль
PATCH /me/profile	Обновить профиль
GET /me/donates	История донатов пользователя
GET /me/summary	Сводка (сумма донатов и т.п.)
POST /donate	Новый донат (привязывается к user_id)
*/

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// Группа маршрутов с авторизацией
	auth := app.Group("/", middleware.RequireAuth(db))

auth.Get("/me/profile", GetProfile(db))
	auth.Patch("/me/profile", UpdateProfile(db))
	auth.Get("/me/donates", GetMyDonates(db))
	auth.Get("/me/summary", GetMySummary(db))

	auth.Post("/donate", func(c *fiber.Ctx) error {
		var donation model.Donation
		if err := c.BodyParser(&donation); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}

		userID := c.Locals("user_id").(string)
		donation.UserID = userID
		donation.CreatedAt = time.Now()

		if err := db.Create(&donation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
		}

		return c.Status(fiber.StatusCreated).JSON(donation)
	})

	// Публичный эндпоинт — список всех донатов
	app.Get("/donates", func(c *fiber.Ctx) error {
		var donations []model.Donation
		if err := db.Order("created_at desc").Find(&donations).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
		}
		return c.JSON(donations)
	})
}
