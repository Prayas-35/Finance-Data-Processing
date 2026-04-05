package handlers

import (
	"context"
	"strconv"

	"github.com/Prayas-35/Finance-Data-Processing/internal/services"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	users *services.UserService
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type activeRequest struct {
	Active bool `json:"active"`
}

func NewUserHandler(users *services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}

	id, err := h.users.Create(context.Background(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"id": id}})
}

func (h *UserHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.users.List(context.Background(), limit, offset)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch users", nil)
	}
	return utils.WriteOK(c, users)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var req updateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}

	if err := h.users.Update(context.Background(), id, req.Name, req.Role); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}
	return utils.WriteOK(c, fiber.Map{"updated": true})
}

func (h *UserHandler) SetActive(c fiber.Ctx) error {
	id := c.Params("id")
	var req activeRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}

	if err := h.users.SetActive(context.Background(), id, req.Active); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}
	return utils.WriteOK(c, fiber.Map{"updated": true})
}
