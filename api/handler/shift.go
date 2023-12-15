package handler

import (
	"context"
	"database/sql"
	"net/http"

	"market_system/config"
	"market_system/pkg/helpers"

	"market_system/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new shift
// @Description Create a new shift in the market system.
// @Tags shift
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param shift body models.CreateShift true "Shift information"
// @Success 201 {object} models.Shift "Created shift"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/shift [post]
func (h *Handler) CreateShift(c *gin.Context) {

	var createShift models.CreateShift
	err := c.ShouldBindJSON(&createShift)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Shift().Create(ctx, &createShift)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a shift by ID
// @Description Get shift details by its ID.
// @Tags shift
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Shift ID"
// @Success 200 {object} models.Shift "Shift details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Shift not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/shift/{id} [get]
func (h *Handler) GetByIDShift(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Shift().GetByID(ctx, &models.ShiftPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusBadRequest, "no rows in result set")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of shifts
// @Description Get a list of shifts with optional filtering.
// @Tags shift
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Shift "List of shifts"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/shift [get]
func (h *Handler) GetListShift(c *gin.Context) {

	limit, err := getIntegerOrDefaultValue(c.Query("limit"), 10)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query limit")
		return
	}

	offset, err := getIntegerOrDefaultValue(c.Query("offset"), 0)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query offset")
		return
	}

	search := c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Shift().GetList(ctx, &models.GetListShiftRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Update a shift
// @Description Update an existing shift.
// @Tags shift
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Shift ID"
// @Param shift body models.UpdateShift true "Updated shift information"
// @Success 202 {object} models.Shift "Updated shift"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Shift not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/shift/{id} [put]
func (h *Handler) UpdateShift(c *gin.Context) {

	var updateShift models.UpdateShift

	err := c.ShouldBindJSON(&updateShift)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}
	updateShift.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Shift().Update(ctx, &updateShift)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(c, http.StatusBadRequest, "no rows affected")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Shift().GetByID(ctx, &models.ShiftPrimaryKey{Id: updateShift.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a shift
// @Description Delete an existing shift.
// @Tags shift
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Shift ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/shift/{id} [delete]
func (h *Handler) DeleteShift(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Shift().Delete(ctx, &models.ShiftPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
