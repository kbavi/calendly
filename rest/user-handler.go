package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/pkg/user"
	"gorm.io/gorm"
)

type UserHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) Create(c *gin.Context) {
	var req pkg.CreateUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDTO := &pkg.UserDTO{}
	userDTO.FromUser(user)

	c.JSON(http.StatusCreated, pkg.ReturnUserResponse{
		Status: "success",
		Data: map[string]pkg.UserDTO{
			"user": *userDTO,
		},
	})
}

func (h *userHandler) Get(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	user, err := h.userService.Get(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDTO := &pkg.UserDTO{}
	userDTO.FromUser(user)

	c.JSON(http.StatusOK, pkg.ReturnUserResponse{
		Status: "success",
		Data:   map[string]pkg.UserDTO{"user": *userDTO},
	})
}

func (h *userHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DeleteUserResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
