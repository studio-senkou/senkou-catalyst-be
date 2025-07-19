package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"
	"strconv"

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
			return throw.ValidationError(c, "Bad request", vErr.Errors)
		}

		return throw.InternalError(c, "Internal server error", fmt.Sprintf("Could not process your request due to an error: %v", err.Error()))
	}

	newUser, appError := h.service.Create(models.User{
		Name:     registerUserDTO.Name,
		Email:    registerUserDTO.Email,
		Password: []byte(registerUserDTO.Password),
	})

	if appError != nil {
		return throw.InternalError(c, "Failed to create user", appError.Details)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data": fiber.Map{
			"user": newUser,
		},
	})
}

// @Summary Get all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserController) GetUsers(c *fiber.Ctx) error {
	users, appError := h.service.GetAll()

	if appError != nil {
		return throw.InternalError(c, "Failed to retrieve users", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data": fiber.Map{
			"users": users,
		},
	})
}

// Get user detail
// @Summary Get user detail
// @Description Get details of a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fiber.Map{data=fiber.Map{user=models.User}}
// @Failure 400 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string, error=string}
// @Router /users/me [get]
func (h *UserController) GetUserDetail(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if userID == 0 || err != nil {
		return throw.Unauthorized(c, "You must be logged in to access this resource")
	}

	user, appError := h.service.GetUserDetail(uint32(userID))

	if appError != nil {
		return throw.InternalError(c, "Failed to retrieve user details", fmt.Sprintf("Could not process your request due to an error: %v", appError.Details))
	}

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "User detail endpoint not implemented",
		"data": fiber.Map{
			"user": user,
		},
	})
}
