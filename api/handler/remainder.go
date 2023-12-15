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

// @Summary Create a new remainder
// @Description Create a new remainder in the market system.
// @Tags remainder
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param remainder body models.CreateRemainder true "Remainder information"
// @Success 201 {object} models.Remainder "Created remainder"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/remainder [post]
func (h *Handler) CreateRemainder(c *gin.Context) {
	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createRemainder models.CreateRemainder
	err := c.ShouldBindJSON(&createRemainder)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON error: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Remainder().Create(ctx, &createRemainder)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a remainder by ID
// @Description Get remainder details by its ID.
// @Tags remainder
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Remainder ID"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/remainder/{id} [get]
func (h *Handler) GetByIDRemainder(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Remainder().GetByID(ctx, &models.RemainderPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusNotFound, "No rows in result set")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of remainder
// @Description Get a list of remainder with optional filtering.
// @Tags remainder
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Remainder "List of remainder"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/remainder [get]
func (h *Handler) GetListRemainder(c *gin.Context) {
	limit, err := getIntegerOrDefaultValue(c.Query("limit"), 10)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Invalid query limit")
		return
	}

	offset, err := getIntegerOrDefaultValue(c.Query("offset"), 0)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Invalid query offset")
		return
	}

	search := c.Query("search")
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Invalid query search")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		key  = fmt.Sprintf("remainder-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListRemainderResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.Remainder) <= 0 {
		resp, err = h.strg.Remainder().GetList(ctx, &models.GetListRemainderRequest{
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

// @Summary Update a remainder
// @Description Update an existing remainder.
// @Tags remainder
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Remainder ID"
// @Param remainder body models.UpdateRemainder true "Updated remainder information"
// @Success 202 {object} models.Remainder "Updated remainder"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/remainder/{id} [put]
func (h *Handler) UpdateRemainder(c *gin.Context) {
	var updateRemainder models.UpdateRemainder

	err := c.ShouldBindJSON(&updateRemainder)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Remainder().Update(ctx, &updateRemainder)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(c, http.StatusBadRequest, "No rows affected")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Remainder().GetByID(ctx, &models.RemainderPrimaryKey{Id: updateRemainder.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a remainder
// @Description Delete an existing remainder.
// @Tags remainder
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Remainder ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/remainder/{id} [delete]
func (h *Handler) DeleteRemainder(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Remainder().Delete(ctx, &models.RemainderPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
