package payment

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

func (h *Handler) CreatePayment(c *gin.Context) {
	var req struct {
		OrderID string  `json:"order_id"`
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pay, err := h.service.CreatePayment(req.OrderID, req.UserID, req.Amount, req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pay)
}

func (h *Handler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	pay, err := h.service.GetPayment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, pay)
}

func (h *Handler) ConfirmPayment(c *gin.Context) {
	id := c.Param("id")
	pay, err := h.service.ConfirmPayment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, pay)
}

func (h *Handler) FailPayment(c *gin.Context) {
	id := c.Param("id")
	pay, err := h.service.FailPayment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, pay)
}

func (h *Handler) GetPaymentsByUser(c *gin.Context) {
	userID := c.Param("id")
	pays, err := h.service.GetPaymentsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pays)
}
