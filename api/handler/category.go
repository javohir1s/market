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

// @Summary Create a new category
// @Description Create a new category in the market system.
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param category body models.CreateCategory true "Category information"
// @Success 201 {object} models.Category "Created category"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/category [post]
func (h *Handler) CreateCategory(c *gin.Context) {

	var createCategory models.CreateCategory
	err := c.ShouldBindJSON(&createCategory)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	if createCategory.ParentID != "" {
		if !helpers.IsValidUUID(createCategory.ParentID) {
			handleResponse(c, http.StatusBadRequest, "parent id is not uuid")
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Category().Create(ctx, &createCategory)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a category by ID
// @Description Get category details by its ID.
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Success 200 {object} models.Category "Category details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/category/{id} [get]
func (h *Handler) GetByIDCategory(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Category().GetByID(ctx, &models.CategoryPrimaryKey{Id: id})
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

// @Summary Get a list of category
// @Description Get a list of category with optional filtering.
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Category "List of category"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/category [get]
func (h *Handler) GetListCategory(c *gin.Context) {

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
		key  = fmt.Sprintf("category-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListCategoryResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.Categories) <= 0 {
		resp, err = h.strg.Category().GetList(ctx, &models.GetListCategoryRequest{
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

// @Summary Update a category
// @Description Update an existing category.
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param Authorization header string true "Authentication token"
// @Param id path string true "Category ID"
// @Param category body models.UpdateCategory true "Updated category information"
// @Success 202 {object} models.Category "Updated category"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/category/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {

	var updateCategory models.UpdateCategory

	err := c.ShouldBindJSON(&updateCategory)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}
	updateCategory.Id = id

	if updateCategory.ParentID != "" {
		if !helpers.IsValidUUID(updateCategory.ParentID) {
			handleResponse(c, http.StatusBadRequest, "parent id is not uuid")
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Category().Update(ctx, &updateCategory)
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

	resp, err := h.strg.Category().GetByID(ctx, &models.CategoryPrimaryKey{Id: updateCategory.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a category
// @Description Delete an existing category.
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/category/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Category().Delete(ctx, &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
