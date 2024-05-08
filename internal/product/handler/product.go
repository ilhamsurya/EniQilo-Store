package handler

import (
	"fmt"
	"net/http"
	"projectsphere/eniqlo-store/internal/product/entity"
	svc "projectsphere/eniqlo-store/internal/product/service"
	"projectsphere/eniqlo-store/pkg/middleware/auth"
	"projectsphere/eniqlo-store/pkg/protocol/msg"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc svc.ProductService
}

func NewProductHandler(productSvc svc.ProductService) ProductHandler {
	return ProductHandler{
		productSvc: productSvc,
	}
}
func (h ProductHandler) Create(c *gin.Context) {

	if c.GetHeader("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, msg.Unauthorization("No authorization header provided"))
		return
	}

	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest("Request body is empty"))
		return
	}
	payload := new(entity.Product)
	err := c.ShouldBindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest(err.Error()))
		return
	}

	userID, err := auth.GetUserIdInsideCtx(c)
	if err != nil {
		fmt.Println(err)
	}

	if containsNull(payload) {
		c.JSON(http.StatusBadRequest, msg.BadRequest("JSON payload contains null values"))
		return
	}

	resp, err := h.productSvc.Create(c.Request.Context(), *payload, userID)
	if err != nil {
		respError := msg.UnwrapRespError(err)
		c.JSON(respError.Code, respError)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Update updates a product.
func (h ProductHandler) Update(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, msg.Unauthorization("No authorization header provided"))
		return
	}

	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest("Request body is empty"))
		return
	}

	payload := new(entity.Product)
	err := c.ShouldBindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest(err.Error()))
		return
	}

	if containsNull(payload) {
		c.JSON(http.StatusBadRequest, msg.BadRequest("JSON payload contains null values"))
		return
	}

	err = h.productSvc.Update(c.Request.Context(), *payload)
	if err != nil {
		respError := msg.UnwrapRespError(err)
		c.JSON(respError.Code, respError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// Delete deletes a product.
func (h ProductHandler) Delete(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, msg.Unauthorization("No authorization header provided"))
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, msg.BadRequest("Product ID is missing"))
		return
	}

	userID, err := auth.GetUserIdInsideCtx(c)
	if err != nil {
		fmt.Println(err)
	}

	err = h.productSvc.Delete(c.Request.Context(), productID, userID)
	if err != nil {
		respError := msg.UnwrapRespError(err)
		c.JSON(respError.Code, respError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func containsNull(param *entity.Product) bool {
	if param == nil {
		return false
	}

	// Check each field for null
	if param.Name == "" || param.SKU == "" || param.Category == "" || param.ImageURL == "" || param.Notes == "" ||
		param.Price <= 0 || param.Stock < 0 || param.Location == "" {
		return true
	}
	return false
}
