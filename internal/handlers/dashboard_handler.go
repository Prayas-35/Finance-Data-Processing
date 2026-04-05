package handlers

import (
	"context"
	"strconv"

	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
	"github.com/Prayas-35/Finance-Data-Processing/internal/services"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type DashboardHandler struct {
	dash *services.DashboardService
}

func NewDashboardHandler(dash *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dash: dash}
}

func (h *DashboardHandler) Summary(c fiber.Ctx) error {
	filter, err := dashboardFilterFromRequest(c)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	data, err := h.dash.Summary(context.Background(), filter)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch summary", nil)
	}
	return utils.WriteOK(c, data)
}

func (h *DashboardHandler) Categories(c fiber.Ctx) error {
	filter, err := dashboardFilterFromRequest(c)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	data, err := h.dash.CategoryTotals(context.Background(), filter)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch category totals", nil)
	}
	return utils.WriteOK(c, data)
}

func (h *DashboardHandler) Trends(c fiber.Ctx) error {
	filter, err := dashboardFilterFromRequest(c)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	data, err := h.dash.TrendsMonthly(context.Background(), filter)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch trends", nil)
	}
	return utils.WriteOK(c, data)
}

func (h *DashboardHandler) Recent(c fiber.Ctx) error {
	filter, err := dashboardFilterFromRequest(c)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	data, err := h.dash.RecentTransactions(context.Background(), filter, limit)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch recent transactions", nil)
	}
	return utils.WriteOK(c, data)
}

func dashboardFilterFromRequest(c fiber.Ctx) (repositories.DashboardFilter, error) {
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)
	fromDate, err := parseOptionalDate(c.Query("from"))
	if err != nil {
		return repositories.DashboardFilter{}, err
	}
	toDate, err := parseOptionalDate(c.Query("to"))
	if err != nil {
		return repositories.DashboardFilter{}, err
	}

	return repositories.DashboardFilter{
		UserID:   userID,
		Role:     role,
		FromDate: fromDate,
		ToDate:   toDate,
	}, nil
}
