package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"recipes/internal/db/models"
	"recipes/internal/handler/error_handler"
)

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the given details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User Data"
// @Success      200 {object} SuccessResponse{data=models.UserResponse}
// @Failure      400 {object} ErrorResponse "Invalid input format"
// @Failure      403 {object} ErrorResponse "Username or email already exists"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var input models.User

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid input format", err)
		return
	}

	// Валидация входных данных
	if err := input.Validate(); err != nil {
		validationMessage := error_handler.ParseValidationErrors(err)
		NewErrorResponse(c, http.StatusBadRequest, validationMessage, err)
		return
	}

	// Создаем пользователя
	userID, err := h.services.CreateUser(c.Request.Context(), &input)
	if err != nil {
		code, constraint := error_handler.ErrorCode(err)
		if code == UniqueViolation {
			if constraint == "users_username_key" {
				NewErrorResponse(c, http.StatusForbidden, "Username is already exist", err)
			}
			if constraint == "users_email_key" {
				NewErrorResponse(c, http.StatusForbidden, "Email is already exist", err)
			}
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err)
		}
		return
	}

	userResponse := models.UserResponse{ID: userID}

	// Формируем ответ
	response := SuccessResponse{
		Status: StatusOK,
		Data:   userResponse,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve a user by their ID
// @Tags         users
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} SuccessResponse{data=models.UserResponse}
// @Failure      400 {object} ErrorResponse "Invalid user ID format"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	user, err := h.services.GetUserById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			NewErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Формируем ответ
	response := SuccessResponse{
		Status: StatusOK,
		Data:   userResponse,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user details by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path int         true "User ID"
// @Param        user body models.UserUpdate true "Updated User Data"
// @Success      200 {object} SuccessResponse{data=string} "User updated successfully"
// @Failure      400 {object} ErrorResponse "Invalid input format"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	var input models.UserUpdate

	// Получаем ID из параметров
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}
	input.ID = id

	// Привязываем входные данные
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid input format", err)
		return
	}

	// Валидируем данные
	if err := input.Validate(); err != nil {
		validationMessage := error_handler.ParseValidationErrors(err)
		NewErrorResponse(c, http.StatusBadRequest, validationMessage, err)
		return
	}

	// Обновляем пользователя
	if err := h.services.UpdateUser(c.Request.Context(), &input); err != nil {
		if err.Error() == "user not found" {
			NewErrorResponse(c, http.StatusNotFound, "User not found", err)
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		}
		return
	}

	response := SuccessResponse{
		Status: StatusOK,
		Data:   "User updated successfully",
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         users
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} SuccessResponse{data=string} "User deleted successfully"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	if err := h.services.DeleteUser(c, id); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Something wrong", err)
		return
	}

	response := SuccessResponse{
		Status: StatusOK,
		Data:   "User deleted successfully",
	}

	c.JSON(http.StatusOK, response)
}

// ListUser godoc
// @Summary      List users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Success      200 {object} SuccessResponse{data=[]models.UserResponse}
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /users [get]
func (h *Handler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.services.ListUser(ctx)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Something wrong", err)
		return
	}

	response := SuccessResponse{
		Status: StatusOK,
		Data:   users,
	}

	c.JSON(http.StatusOK, response)
}
