package controllers

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

// @Summary Create user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Router /users [post]
func (h *UserController) CreateUser(c *fiber.Ctx) error {
	registerUserDTO := new(dtos.RegisterUserDTO)

	if err := utils.Validate(c, registerUserDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	newUser, err := h.service.Create(models.User{
		Name:     registerUserDTO.Name,
		Email:    registerUserDTO.Email,
		Password: []byte(registerUserDTO.Password),
	})

	if err != nil {
		return throw.InternalError(c, "Failed to create user", map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(newUser)
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
