package handler

import (
	"net/http"
	"strconv"

	"github.com/ZXstrike/marketplace-app/internal/domain/api/service"
	"github.com/ZXstrike/shared/pkg/models"
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
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description" binding:"required"`
		BaseURL      string   `json:"base_url" binding:"required"`
		PricePerCall float64  `json:"price_per_call" binding:"required"`
		Categories   []string `json:"categories"` // Optional, can be empty
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

	err := h.service.CreateNewAPI(req.Name, req.Description, userId.(string), req.BaseURL, req.PricePerCall, req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "API created successfully"})
}

func (h *Handler) UpdateAPI(c *gin.Context) {
	var req struct {
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description" binding:"required"`
		BaseURL      string   `json:"base_url" binding:"required"`
		PricePerCall float64  `json:"price_per_call" binding:"required"`
		Categories   []string `json:"categories"` // Optional, can be empty
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiID := c.Param("id")
	if apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	}

	err := h.service.UpdateAPI(apiID, req.Name, req.Description, req.BaseURL, req.PricePerCall, req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API updated successfully"})
}

func (h *Handler) DeleteAPI(c *gin.Context) {
	apiID := c.Param("id")
	if apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.service.DeleteAPI(userId.(string), apiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API deleted successfully"})
}

func (h *Handler) GetAllAPIs(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}

	length := c.Query("length")
	if length == "" {
		length = "10" // Default to 10 items per page if not provided
	}
	// Convert page and length to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid length"})
		return
	}

	apis, err := h.service.GetAllAPIs(pageInt, lengthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch APIs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, apis)
}

func (h *Handler) GetAPIByID(c *gin.Context) {
	if apiID := c.Param("id"); apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	} else {
		api, err := h.service.GetAPIByID(apiID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API: " + err.Error()})
			return
		}
		if api == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}
		c.JSON(http.StatusOK, api)
	}
}

func (h *Handler) CreateNewAPIEndpoint(c *gin.Context) {
	type baseEndpoint struct {
		HTTPMethod    string `json:"http_method" binding:"required"`
		Path          string `json:"path" binding:"required"`
		Documentation string `json:"documentation" binding:"required"`
	}

	var req struct {
		APIVersionID string         `json:"api_version_id" binding:"required"`
		Endpoints    []baseEndpoint `json:"endpoints" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var endpoints []models.Endpoint
	for _, ep := range req.Endpoints {
		endpoints = append(
			endpoints,
			models.Endpoint{
				HTTPMethod:    ep.HTTPMethod,
				Path:          ep.Path,
				Documentation: ep.Documentation,
			},
		)
	}

	err := h.service.CreateAPIEndpoint(req.APIVersionID, endpoints)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API endpoint: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "API endpoint created successfully"})
}

func (h *Handler) UpdateAPIEndpoint(c *gin.Context) {

	type baseEndpoint struct {
		EndpointID    string `json:"endpoint_id" binding:"required"`
		HTTPMethod    string `json:"http_method" binding:"required"`
		Path          string `json:"path" binding:"required"`
		Documentation string `json:"documentation" binding:"required"`
	}

	var req struct {
		APIVersionID string         `json:"api_version_id" binding:"required"`
		Endpoints    []baseEndpoint `json:"endpoints" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var endpoints []models.Endpoint
	for _, ep := range req.Endpoints {
		endpoints = append(
			endpoints,
			models.Endpoint{
				ID:            ep.EndpointID,
				HTTPMethod:    ep.HTTPMethod,
				Path:          ep.Path,
				Documentation: ep.Documentation,
			},
		)
	}

	err := h.service.UpdateAPIEndpoint(req.APIVersionID, endpoints)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API endpoint: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API endpoint updated successfully"})
}

func (h *Handler) DeleteAPIEndpoint(c *gin.Context) {
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint ID is required"})
		return
	}

	err := h.service.DeleteAPIEndpoint(endpointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API endpoint: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API endpoint deleted successfully"})
}

func (h *Handler) GetAllAPIEndpointsByAPIVersionID(c *gin.Context) {
	apiVersionID := c.Param("apiVersionID")
	if apiVersionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Version ID is required"})
		return
	}

	endpoints, err := h.service.GetAllEndpointsByAPIVersionID(apiVersionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API endpoints: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, endpoints)
}
