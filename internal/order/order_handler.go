package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerName  string      `json:"customer_name"`
		Phone         string      `json:"phone"`
		Email         string      `json:"email"`
		Address       string      `json:"address"`
		PaymentMethod string      `json:"payment_method"`
		FinalTotal    float64     `json:"final_total"`
		Items         []OrderItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := h.service.CreateOrder(
		req.CustomerName,
		req.Phone,
		req.Email,
		req.Address,
		req.PaymentMethod,
		req.FinalTotal,
		req.Items,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	ord, err := h.service.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, ord)
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.service.UpdateOrderStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteOrder(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("id")
	orders, err := h.service.GetOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) ListOrders(c *gin.Context) {
	all, err := h.service.ListAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, all)
}
