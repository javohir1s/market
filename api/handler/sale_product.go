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

// @Summary Create a new sale product
// @Description Create a new sale product in the market system.
// @Tags sale_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param saleProduct body models.CreateSaleProduct true "Sale Product information"
// @Success 201 {object} models.SaleProduct "Created sale product"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_product [post]
func (h *Handler) CreateSaleProduct(c *gin.Context) {
	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createSaleProduct models.CreateSaleProduct
	err := c.ShouldBindJSON(&createSaleProduct)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Product().Create(ctx, &createSaleProduct)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a sale product by ID
// @Description Get sale product details by its ID.
// @Tags sale_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Product ID"
// @Success 200 {object} models.SaleProduct "Sale Product details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Sale Product not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_product/{id} [get]
func (h *Handler) GetByIDSaleProduct(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Product().GetByID(ctx, &models.SaleProductPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusNotFound, "sale product not found")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of sale products
// @Description Get a list of sale products with optional filtering.
// @Tags sale_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {object} models.GetListSaleProductResponse "List of sale products"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_product [get]
func (h *Handler) GetListSaleProduct(c *gin.Context) {
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

	resp, err := h.strg.Sale_Product().GetList(ctx, &models.GetListSaleProductRequest{
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

// @Summary Update a sale product
// @Description Update an existing sale product.
// @Tags sale_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Product ID"
// @Param saleProduct body models.UpdateSaleProduct true "Updated sale product information"
// @Success 202 {object} models.SaleProduct "Updated sale product"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Sale Product not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_product/{id} [put]
func (h *Handler) UpdateSaleProduct(c *gin.Context) {
	var updateSaleProduct models.UpdateSaleProduct

	err := c.ShouldBindJSON(&updateSaleProduct)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}
	updateSaleProduct.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Sale_Product().Update(ctx, &updateSaleProduct)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(c, http.StatusNotFound, "sale product not found")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale_Product().GetByID(ctx, &models.SaleProductPrimaryKey{Id: updateSaleProduct.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a sale product
// @Description Delete an existing sale product.
// @Tags sale_product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Sale Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/sale_product/{id} [delete]
func (h *Handler) DeleteSaleProduct(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Sale_Product().Delete(ctx, &models.SaleProductPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
