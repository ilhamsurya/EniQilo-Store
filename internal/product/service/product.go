package svc

import (
	"context"
	"fmt"
	"net/http"
	"projectsphere/eniqlo-store/internal/product/entity"
	"projectsphere/eniqlo-store/internal/product/repository"
	"projectsphere/eniqlo-store/pkg/protocol/msg"
	"strings"
)

type ProductService struct {
	productRepo repository.ProductRepo
}

func NewProductService(productRepo repository.ProductRepo) ProductService {
	return ProductService{
		productRepo: productRepo,
	}
}

func (s ProductService) Update(ctx context.Context, product entity.Product) error {
	if err := s.validateProduct(product); err != nil {
		return &msg.RespError{
			Code:    http.StatusBadRequest,
			Message: "request doesn't pass validation",
		}
	}

	err := s.productRepo.UpdateProduct(product)
	if err != nil {
		return err
	}

	return nil
}

func (s ProductService) Delete(ctx context.Context, productID string, userID uint32) error {
	err := s.productRepo.DeleteProduct(productID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s ProductService) Create(ctx context.Context, productParam entity.Product, userId uint32) (entity.ProductResponse, error) {
	if err := s.validateProduct(productParam); err != nil {
		return entity.ProductResponse{}, &msg.RespError{
			Code:    http.StatusBadRequest,
			Message: "request doesn't pass validation",
		}
	}
	product, err := s.productRepo.CreateProduct(ctx, productParam, userId)
	if err != nil {
		return entity.ProductResponse{}, err
	}

	return entity.ProductResponse{
		Message: "success",
		Data: entity.Product{
			ID:        product.ID,
			CreatedAt: product.CreatedAt,
		},
	}, nil
}

func (s ProductService) validateProduct(product entity.Product) error {
	validationErrors := make(map[string]string)

	if product.Name == "" {
		validationErrors["name"] = "name cannot be empty"
	}
	if product.SKU == "" {
		validationErrors["sku"] = "sku cannot be empty"
	}
	if product.Category == "" {
		validationErrors["category"] = "category cannot be empty"
	}
	if product.ImageURL == "" {
		validationErrors["imageUrl"] = "imageUrl cannot be empty"
	}
	if product.Notes == "" {
		validationErrors["notes"] = "notes cannot be empty"
	}
	if product.Price == 0 {
		validationErrors["price"] = "price cannot be zero"
	}
	if product.Stock == 0 {
		validationErrors["stock"] = "stock cannot be zero"
	}
	if product.Location == "" {
		validationErrors["location"] = "location cannot be empty"
	}

	if len(product.Name) < 1 || len(product.Name) > 30 {
		validationErrors["name"] = "name must be between 1 and 30 characters"
	}

	if len(product.SKU) < 1 || len(product.SKU) > 30 {
		validationErrors["sku"] = "sku must be between 1 and 30 characters"
	}

	validCategories := map[string]bool{
		"Clothing":    true,
		"Accessories": true,
		"Footwear":    true,
		"Beverages":   true,
	}
	if !validCategories[product.Category] {
		validationErrors["category"] = "invalid category"
	}

	if !strings.HasPrefix(product.ImageURL, "http://") && !strings.HasPrefix(product.ImageURL, "https://") {
		validationErrors["imageUrl"] = "invalid URL"
	}

	if len(product.Notes) < 1 || len(product.Notes) > 200 {
		validationErrors["notes"] = "notes must be between 1 and 200 characters"
	}

	if product.Price < 1 {
		validationErrors["price"] = "price must be greater than 0"
	}

	if product.Stock < 0 || product.Stock > 100000 {
		validationErrors["stock"] = "stock must be between 0 and 100000"
	}

	if len(product.Location) < 1 || len(product.Location) > 200 {
		validationErrors["location"] = "location must be between 1 and 200 characters"
	}

	if len(validationErrors) > 0 {
		var errorMsgs []string
		for field, msg := range validationErrors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", field, msg))
		}
		return fmt.Errorf(strings.Join(errorMsgs, "; "))
	}

	return nil
}

func (s ProductService) Get(ctx context.Context, getProductParam entity.GetProductParam) ([]entity.GetProductData, error) {
	orderby := []string{"asc", "desc"}
	isValidPricePrm := false
	isValidCreatePrm := false

	for _, v := range orderby {
		if getProductParam.Price == v {
			isValidPricePrm = true
		}

		if getProductParam.CreatedAt == v {
			isValidCreatePrm = true
		}
	}

	if !isValidCreatePrm {
		getProductParam.CreatedAt = ""
	}

	if !isValidPricePrm {
		getProductParam.Price = ""
	}

	res, err := s.productRepo.GetProduct(ctx, getProductParam)
	if err != nil {
		return nil, err
	}

	return res, nil
}
