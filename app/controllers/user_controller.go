package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/query"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService services.UserService
	subService  services.SubscriptionService
}

func NewUserController(userService services.UserService, subService services.SubscriptionService) *UserController {
	return &UserController{
		userService: userService,
		subService:  subService,
	}
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

	if err := validator.Validate(c, registerUserDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Bad request", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", fmt.Sprintf("Could not process your request due to an error: %v", err.Error()))
	}

	newUser, appError := h.userService.Create(models.User{
		Name:     registerUserDTO.Name,
		Email:    registerUserDTO.Email,
		Phone:    registerUserDTO.Phone,
		Password: []byte(registerUserDTO.Password),
	})

	if appError != nil {
		switch appError.Code {
		case fiber.StatusBadRequest:
			return response.BadRequest(c, "Cannot continue to register user, user already exists", appError.Details)
		default:
			return response.InternalError(c, "Failed to create user", appError.Details)
		}
	}

	if err := h.subService.AssignFreeTierSubscription(newUser.ID); err != nil {
		return response.InternalError(c, "Failed to assign free tier subscription", err.Details)
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
	params := query.ParseQueryParams(c)

	users, pagination, appError := h.userService.GetAll(params)

	if appError != nil {
		return response.InternalError(c, "Failed to retrieve users", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data": fiber.Map{
			"users":      users,
			"pagination": pagination,
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
		return response.Unauthorized(c, "You must be logged in to access this resource")
	}

	user, appError := h.userService.GetUserDetail(uint32(userID))

	if appError != nil {
		return response.InternalError(c, "Failed to retrieve user details", fmt.Sprintf("Could not process your request due to an error: %v", appError.Details))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User detail endpoint not implemented",
		"data": fiber.Map{
			"user": user,
		},
	})
}
