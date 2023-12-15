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

// @Summary Create a new branch
// @Description Create a new branch in the market system.
// @Tags branch
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param branch body models.CreateBranch true "Branch information"
// @Success 201 {object} models.Branch "Created branch"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/branch [post]
func (h *Handler) Createbranch(c *gin.Context) {

	var createbranch models.CreateBranch
	err := c.ShouldBindJSON(&createbranch)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Branch().Create(ctx, &createbranch)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a branch by ID
// @Description Get branch details by its ID.
// @Tags branch
// @Accept json
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/branch/{id} [get]
func (h *Handler) GetByIDbranch(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Branch().GetByID(ctx, &models.BranchPrimaryKey{Id: id})
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

// @Summary Get a list of branch
// @Description Get a list of branch with optional filtering.
// @Tags branch
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Branch "List of branch"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/branch [get]
func (h *Handler) GetListbranch(c *gin.Context) {

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
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		key  = fmt.Sprintf("branch-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListBranchResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.Branches) <= 0 {
		resp, err = h.strg.Branch().GetList(ctx, &models.GetListBranchRequest{
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

// @Summary Update a branch
// @Description Update an existing branch.
// @Tags branch
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param id path string true "Branch ID"
// @Param branch body models.UpdateBranch true "Updated branch information"
// @Success 202 {object} models.Branch "Updated branch"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/branch/{id} [put]
func (h *Handler) Updatebranch(c *gin.Context) {

	var updatebranch models.UpdateBranch

	err := c.ShouldBindJSON(&updatebranch)
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

	rowsAffected, err := h.strg.Branch().Update(ctx, &updatebranch)
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

	resp, err := h.strg.Branch().GetByID(ctx, &models.BranchPrimaryKey{Id: updatebranch.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a branch
// @Description Delete an existing branch.
// @Tags branch
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param id path string true "Branch ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/branch/{id} [delete]
func (h *Handler) Deletebranch(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Branch().Delete(ctx, &models.BranchPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
