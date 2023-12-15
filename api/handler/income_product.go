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

// @Summary Create a new income product
// @Description Create a new income product in the market system.
// @Tags income_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param income_product body models.CreateIncomeProduct true "Income product information"
// @Success 201 {object} models.IncomeProduct "Created income product"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income_product [post]
func (h *Handler) CreateIncomeProduct(c *gin.Context) {
	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createIncomeProduct models.CreateIncomeProduct
	err := c.ShouldBindJSON(&createIncomeProduct)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON error: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.IncomeProduct().Create(ctx, &createIncomeProduct)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get an income product by ID
// @Description Get income product details by its ID.
// @Tags income_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income Product ID"
// @Success 200 {object} models.IncomeProduct "Income product details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Income product not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income_product/{id} [get]
func (h *Handler) GetByIDIncomeProduct(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.IncomeProduct().GetByID(ctx, &models.IncomeProductPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusBadRequest, "No rows in the result set")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of income products
// @Description Get a list of income products with optional filtering.
// @Tags income_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.IncomeProduct "List of income products"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income_product [get]
func (h *Handler) GetListIncomeProduct(c *gin.Context) {
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
		key  = fmt.Sprintf("income_product-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListIncomeProductResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.IncomeProducts) <= 0 {
		resp, err = h.strg.IncomeProduct().GetList(ctx, &models.GetListIncomeProductRequest{
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

// @Summary Update an income product
// @Description Update an existing income product.
// @Tags income_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income Product ID"
// @Param income_product body models.UpdateIncomeProduct true "Updated income product information"
// @Success 202 {object} models.IncomeProduct "Updated income product"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Income product not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income_product/{id} [put]
func (h *Handler) UpdateIncomeProduct(c *gin.Context) {
	var updateIncomeProduct models.UpdateIncomeProduct
	err := c.ShouldBindJSON(&updateIncomeProduct)
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

	rowsAffected, err := h.strg.IncomeProduct().Update(ctx, &updateIncomeProduct)
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

	resp, err := h.strg.IncomeProduct().GetByID(ctx, &models.IncomeProductPrimaryKey{Id: updateIncomeProduct.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete an income product
// @Description Delete an existing income product.
// @Tags income_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Income Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/income_product/{id} [delete]
func (h *Handler) DeleteIncomeProduct(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.IncomeProduct().Delete(ctx, &models.IncomeProductPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
