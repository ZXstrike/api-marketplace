package handler

import (
	"net/http"

	"github.com/ZXstrike/marketplace-app/internal/domain/api/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateNewAPI(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.service.CreateNewAPI(req.Name, req.Description, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "API created successfully"})
}

// func (h *Handler) GetAllAPIs(c *gin.Context) {
// 	return h.service.GetAllAPIs()
// }

// func (h *Handler) GetAPIByID(c *gin.Context) {
// 	return h.service.GetAPIByID(id)
// }

// func (h *Handler) GetAllAPIsByUserID(c *gin.Context) {
// 	return h.service.GetAllAPIsByUserID(userID)
// }

// func (h *Handler) UpdateAPI(c *gin.Context) {
// 	return h.service.UpdateAPI(api)
// }

// func (h *Handler) DeleteAPI(c *gin.Context) {
// 	return h.service.DeleteAPI(userId, apiId)
// }
