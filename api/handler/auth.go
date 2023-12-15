package handler

import (
	"net/http"

	"market_system/config"
	"market_system/models"
	"market_system/pkg/security"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @Summary User Login
// @Description Logs in a user and returns an access token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "User login information"
// @Success 200 {object} models.LoginResponse "Successful login"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {

	var req models.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "ShouldBindJSON err:"+err.Error())
		return
	}

	user, err := h.strg.User().GetByID(c.Request.Context(), &models.UserPrimaryKey{Id: req.Login})
	if err != nil {
		if err == pgx.ErrNoRows {
			handleResponse(c, http.StatusBadRequest, "not found user")
			return
		}

		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user.Password != req.Password {
		handleResponse(c, http.StatusBadRequest, "invalid login or password")
		return
	}

	var credentails = map[string]interface{}{
		"user_id":     user.Id,
		"client_type": user.ClientType,
	}
	accessToken, err := security.GenerateJWT(credentails, config.ExpiredTime, h.cfg.SecretKey)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		AccessToken: accessToken,
		User:        *user,
	}

	handleResponse(c, http.StatusOK, resp)
}
