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

// @Summary Create a new transaction
// @Description Create a new transaction in the market system.
// @Tags transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param transaction body models.CreateTransaction true "Transaction information"
// @Success 201 {object} models.Transaction "Created transaction"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/transaction [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createTransaction models.CreateTransaction
	err := c.ShouldBindJSON(&createTransaction)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Transaction().Create(ctx, &createTransaction)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a transaction by ID
// @Description Get transaction details by its ID.
// @Tags transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.Transaction "Transaction details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Transaction not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/transaction/{id} [get]
func (h *Handler) GetByIDTransaction(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Transaction().GetByID(ctx, &models.TransactionPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusNotFound, "no rows in result set")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of transactions
// @Description Get a list of transactions with optional filtering.
// @Tags transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Transaction "List of transactions"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/transaction [get]
func (h *Handler) GetListTransaction(c *gin.Context) {
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

	resp, err := h.strg.Transaction().GetList(ctx, &models.GetListTransactonRequest{
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

// @Summary Update a transaction
// @Description Update an existing transaction.
// @Tags transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Transaction ID"
// @Param transaction body models.UpdateTransaction true "Updated transaction information"
// @Success 202 {object} models.Transaction "Updated transaction"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Transaction not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/transaction/{id} [put]
func (h *Handler) UpdateTransaction(c *gin.Context) {
	var updateTransaction models.UpdateTransaction
	err := c.ShouldBindJSON(&updateTransaction)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}
	updateTransaction.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Transaction().Update(ctx, &updateTransaction)
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

	resp, err := h.strg.Transaction().GetByID(ctx, &models.TransactionPrimaryKey{Id: updateTransaction.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a transaction
// @Description Delete an existing transaction.
// @Tags transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Transaction ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/transaction/{id} [delete]
func (h *Handler) DeleteTransaction(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Transaction().Delete(ctx, &models.TransactionPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
