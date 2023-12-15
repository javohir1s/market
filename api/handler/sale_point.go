package handler

import (
	"context"
	"database/sql"
	"net/http"

	"market_system/config"
	"market_system/models"
	"market_system/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new sale point
// @Description Create a new sale point in the market system.
// @Tags sale_point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param salePoint body models.CreateSalePoint true "Sale point information"
// @Success 201 {object} models.SalePoint "Created sale point"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_point [post]
func (h *Handler) CreateSalePoint(c *gin.Context) {
	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createSalePoint models.CreateSalePoint
	if err := c.ShouldBindJSON(&createSalePoint); err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON error: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Point().Create(ctx, &createSalePoint)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a sale point by ID
// @Description Get sale point details by its ID.
// @Tags sale_point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Point ID"
// @Success 200 {object} models.SalePoint "Sale point details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Sale point not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_point/{id} [get]
func (h *Handler) GetByIDSalePoint(c *gin.Context) {
	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Point().GetByID(ctx, &models.SalePointPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusNotFound, "Sale point not found")
		return
	} else if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of sale points
// @Description Get a list of sale points with optional filtering.
// @Tags sale_point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.SalePoint "List of sale points"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_point [get]
func (h *Handler) GetListSalePoint(c *gin.Context) {
	limit, err := getIntegerOrDefaultValue(c.Query("limit"), 10)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Invalid query parameter 'limit'")
		return
	}

	offset, err := getIntegerOrDefaultValue(c.Query("offset"), 0)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Invalid query parameter 'offset'")
		return
	}

	search := c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Point().GetList(ctx, &models.GetListSalePointRequest{
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

// @Summary Update a sale point
// @Description Update an existing sale point.
// @Tags sale_point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Point ID"
// @Param salePoint body models.UpdateSalePoint true "Updated sale point information"
// @Success 202 {object} models.SalePoint "Updated sale point"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Sale point not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_point/{id} [put]
func (h *Handler) UpdateSalePoint(c *gin.Context) {
	var updateSalePoint models.UpdateSalePoint

	err := c.ShouldBindJSON(&updateSalePoint)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}
	updateSalePoint.Id = id

	if updateSalePoint.Branch_id != "" && !helpers.IsValidUUID(updateSalePoint.Branch_id) {
		handleResponse(c, http.StatusBadRequest, "parent ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Sale_Point().Update(ctx, &updateSalePoint)
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

	resp, err := h.strg.Sale_Point().GetByID(ctx, &models.SalePointPrimaryKey{Id: updateSalePoint.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a sale point
// @Description Delete an existing sale point.
// @Tags sale_point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Point ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_point/{id} [delete]
func (h *Handler) DeleteSalePoint(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Sale_Point().Delete(ctx, &models.SalePointPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
