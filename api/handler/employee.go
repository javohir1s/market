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

// @Summary Create a new employee
// @Description Create a new employee in the market system.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param employee body models.CreateEmployee true "Employee information"
// @Success 201 {object} models.Employee "Created employee"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/employee [post]
func (h *Handler) CreateEmployee(c *gin.Context) {

	password := c.GetHeader("Password")
	if password != "1234" {
		handleResponse(c, http.StatusUnauthorized, "The request requires user authentication.")
		return
	}

	var createEmployee models.CreateEmployee
	err := c.ShouldBindJSON(&createEmployee)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON error: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Employee().Create(ctx, &createEmployee)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get an employee by ID
// @Description Get employee details by its ID.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Employee ID"
// @Success 200 {object} models.Employee "Employee details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Employee not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/employee/{id} [get]
func (h *Handler) GetByIDEmployee(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Employee().GetByID(ctx, &models.EmployeePrimaryKey{Id: id})
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

// @Summary Get a list of employees
// @Description Get a list of employees with optional filtering.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param search query string false "Search term"
// @Success 200 {array} models.Employee "List of employees"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/employee [get]
func (h *Handler) GetListEmployee(c *gin.Context) {

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
		key  = fmt.Sprintf("employee-%s", c.Request.URL.Query().Encode())
		resp = &models.GetListEmployeeResponse{}
	)

	body, err := h.cache.GetX(ctx, key)
	if err == nil {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(resp.Employees) <= 0 {
		resp, err = h.strg.Employee().GetList(ctx, &models.GetListEmployeeRequest{
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

// @Summary Update an employee
// @Description Update an existing employee.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Employee ID"
// @Param employee body models.UpdateEmployee true "Updated employee information"
// @Success 202 {object} models.Employee "Updated employee"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Employee not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/employee/{id} [put]
func (h *Handler) UpdateEmployee(c *gin.Context) {

	var updateEmployee models.UpdateEmployee

	err := c.ShouldBindJSON(&updateEmployee)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}
	updateEmployee.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Employee().Update(ctx, &updateEmployee)
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

	resp, err := h.strg.Employee().GetByID(ctx, &models.EmployeePrimaryKey{Id: updateEmployee.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete an employee
// @Description Delete an existing employee.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication token"
// @Param Password header string true "User password"
// @Param id path string true "Employee ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/employee/{id} [delete]
func (h *Handler) DeleteEmployee(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "ID is not a valid UUID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Employee().Delete(ctx, &models.EmployeePrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
