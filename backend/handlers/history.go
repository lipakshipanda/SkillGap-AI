package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"skillgap-ai/models"
	"skillgap-ai/services"
)

type HistoryHandler struct {
	Store *services.Store
}

func NewHistoryHandler(store *services.Store) *HistoryHandler {
	return &HistoryHandler{Store: store}
}

// GET /api/history?limit=10
func (h *HistoryHandler) GetHistory(c *fiber.Ctx) error {
	limit := 10
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}

	items, err := h.Store.GetHistory(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load history")
	}
	if items == nil {
		items = []models.HistoryItem{} // keep JSON output as [] instead of null
	}
	return c.JSON(items)
}

// GET /api/history/:id
func (h *HistoryHandler) GetAnalysis(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.Store.GetAnalysis(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load analysis")
	}
	if resp == nil {
		return fiber.NewError(fiber.StatusNotFound, "Analysis not found")
	}
	return c.JSON(resp)
}

// DELETE /api/history/:id
func (h *HistoryHandler) DeleteAnalysis(c *fiber.Ctx) error {
	id := c.Params("id")
	deleted, err := h.Store.DeleteAnalysis(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete analysis")
	}
	if !deleted {
		return fiber.NewError(fiber.StatusNotFound, "Analysis not found")
	}
	return c.JSON(fiber.Map{"deleted": true})
}
