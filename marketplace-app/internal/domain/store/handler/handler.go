package handler

import (
	"net/http"

	"github.com/ZXstrike/marketplace-app/internal/domain/store/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetStoreByUserIDHandler(c *gin.Context) {
	userID := c.Param("userID")
	store, err := h.service.GetStoreByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get store"})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (h *Handler) GetStoreByUsernameHandler(c *gin.Context) {
	username := c.Param("username")
	store, err := h.service.GetStoreByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get store"})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (h *Handler) GetAllStoresHandler(c *gin.Context) {
	stores, err := h.service.GetAllStores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stores"})
		return
	}
	c.JSON(http.StatusOK, stores)
}

func (h *Handler) CreateStoreHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.CreateStore(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store created successfully"})
}

func (h *Handler) UpdateStoreHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.UpdateStore(c.Request.Context(), userID.(string), req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store updated successfully"})
}
