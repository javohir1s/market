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

// @Summary Create a new user
// @Description Create a new user in the market system.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User information"
// @Success 201 {object} models.User "Created user"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/user [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var createUser models.CreateUser
	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.User().Create(ctx, &createUser)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a user by ID
// @Description Get user details by user ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/user/{id} [get]
func (h *Handler) GetByIDUser(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.User().GetByID(ctx, &models.UserPrimaryKey{Id: id})
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

// @Summary Get a list of users
// @Description Get a list of users with optional filtering.
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.User "List of users"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/user [get]
func (h *Handler) GetListUser(c *gin.Context) {

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

	resp, err := h.strg.User().GetList(ctx, &models.GetListUserRequest{
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

// @Summary Update a user
// @Description Update an existing user.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUser true "Updated user information"
// @Success 202 {object} models.User "Updated user"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/user/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {

	var updateUser models.UpdateUser

	err := c.ShouldBindJSON(&updateUser)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}
	updateUser.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.User().Update(ctx, &updateUser)
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

	resp, err := h.strg.User().GetByID(ctx, &models.UserPrimaryKey{Id: updateUser.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete a user
// @Description Delete an existing user.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.User().Delete(ctx, &models.UserPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
