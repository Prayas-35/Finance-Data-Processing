package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
	"github.com/Prayas-35/Finance-Data-Processing/internal/services"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type RecordHandler struct {
	records *services.RecordService
}

type recordRequest struct {
	Amount   string `json:"amount"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Notes    string `json:"notes"`
}

func NewRecordHandler(records *services.RecordService) *RecordHandler {
	return &RecordHandler{records: records}
}

func (h *RecordHandler) Create(c fiber.Ctx) error {
	var req recordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}

	userID, _ := c.Locals("user_id").(string)
	parsedDate, err := parseDateOrNow(req.Date)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid date format; use RFC3339", nil)
	}

	id, err := h.records.Create(context.Background(), models.FinancialRecord{
		UserID:   userID,
		Amount:   req.Amount,
		Type:     models.RecordType(req.Type),
		Category: req.Category,
		Date:     parsedDate,
		Notes:    req.Notes,
	})
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"id": id}})
}

func (h *RecordHandler) List(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	fromDate, err := parseOptionalDate(c.Query("from"))
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid from date; use RFC3339", nil)
	}
	toDate, err := parseOptionalDate(c.Query("to"))
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid to date; use RFC3339", nil)
	}

	records, err := h.records.List(context.Background(), repositories.RecordFilter{
		UserID:   userID,
		Role:     role,
		FromDate: fromDate,
		ToDate:   toDate,
		Category: c.Query("category"),
		Type:     c.Query("type"),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "query_failed", "could not fetch records", nil)
	}
	return utils.WriteOK(c, records)
}

func (h *RecordHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)

	record, err := h.records.GetByID(context.Background(), id, userID, role)
	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "not_found", "record not found", nil)
	}
	return utils.WriteOK(c, record)
}

func (h *RecordHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)

	var req recordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}
	parsedDate, err := parseDateOrNow(req.Date)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid date format; use RFC3339", nil)
	}

	err = h.records.Update(context.Background(), id, userID, role, models.FinancialRecord{
		Amount:   req.Amount,
		Type:     models.RecordType(req.Type),
		Category: req.Category,
		Date:     parsedDate,
		Notes:    req.Notes,
	})
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", err.Error(), nil)
	}
	return utils.WriteOK(c, fiber.Map{"updated": true})
}

func (h *RecordHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)

	if err := h.records.SoftDelete(context.Background(), id, userID, role); err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "not_found", "record not found", nil)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseDateOrNow(raw string) (time.Time, error) {
	if raw == "" {
		return time.Now().UTC(), nil
	}
	return time.Parse(time.RFC3339, raw)
}

func parseOptionalDate(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
