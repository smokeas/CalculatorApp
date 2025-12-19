package handlers

import (
	"net/http"

	"CalculatorApp/internal/calculationService"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service calculationService.CalcService
}

func NewHandlers(s calculationService.CalcService) *Handlers {
	return &Handlers{service: s}
}

func (h *Handlers) GetCalculations(c echo.Context) error {
	list, err := h.service.GetAllCalculations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handlers) PostCalculation(c echo.Context) error {
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Expression == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "expression required"})
	}

	// create with a generated ID
	created, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// if created has no ID (repo/db should set), assign one for sqlite fallback
	if created.ID == "" {
		created.ID = uuid.New().String()
	}

	return c.JSON(http.StatusCreated, created)
}

func (h *Handlers) PatchCalculation(c echo.Context) error {
	id := c.Param("id")
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Expression == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "expression required"})
	}
	updated, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *Handlers) DeleteCalculation(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteCalculation(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
