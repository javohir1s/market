package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"market_system/config"
	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new income
// @Description Create a new income in the market system.
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param income body models.CreateIncome true "Income information"
// @Success 201 {object} models.Income "Created income"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income [post]
func (h *Handler) CreateIncome(c *gin.Context) {

	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createIncome models.CreateIncome
	err := c.ShouldBindJSON(&createIncome)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Income().Create(ctx, &createIncome)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get an income by ID
// @Description Get income details by its ID.
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income ID"
// @Success 200 {object} models.Income "Income details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Income not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income/{id} [get]
func (h *Handler) GetByIDIncome(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Income().GetByID(ctx, &models.IncomePrimaryKey{Id: id})
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

// @Summary Get a list of incomes
// @Description Get a list of incomes with optional filtering.
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Income "List of incomes"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income [get]
func (h *Handler) GetListIncome(c *gin.Context) {

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
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query search")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		key  = fmt.Sprintf("income-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListIncomeResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.Incomes) <= 0 {
		resp, err = h.strg.Income().GetList(ctx, &models.GetListIncomeRequest{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

		body, err := json.Marshal(resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

		h.cache.SetX(ctx, key, string(body), time.Second*15)
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Update an income
// @Description Update an existing income.
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income ID"
// @Param income body models.UpdateIncome true "Updated income information"
// @Success 202 {object} models.Income "Updated income"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Income not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income/{id} [put]
func (h *Handler) UpdateIncome(c *gin.Context) {

	var updateIncome models.UpdateIncome

	err := c.ShouldBindJSON(&updateIncome)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Income().Update(ctx, &updateIncome)
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

	resp, err := h.strg.Income().GetByID(ctx, &models.IncomePrimaryKey{Id: updateIncome.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete an income
// @Description Delete an existing income.
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income/{id} [delete]
func (h *Handler) DeleteIncome(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Income().Delete(ctx, &models.IncomePrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
