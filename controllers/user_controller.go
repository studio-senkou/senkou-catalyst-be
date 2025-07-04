package controllers

import (
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

// @Summary Get all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserController) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// @Summary Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Router /users [post]
func (h *UserController) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	newUser, err := h.service.Create(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(newUser)
}
