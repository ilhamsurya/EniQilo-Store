package repository

import (
	"context"
	"database/sql"
	"projectsphere/eniqlo-store/internal/product/entity"
	"projectsphere/eniqlo-store/pkg/database"
	"projectsphere/eniqlo-store/pkg/protocol/msg"

	"github.com/google/uuid"
)

type ProductRepo struct {
	dbConnector database.PostgresConnector
}

func NewProductRepo(dbConnector database.PostgresConnector) ProductRepo {
	return ProductRepo{
		dbConnector: dbConnector,
	}
}

// DeleteProduct implements ProductInterface.
func (r ProductRepo) DeleteProduct(id string, userId uuid.UUID) error {
	panic("unimplemented")
}

// UpdateProduct implements ProductInterface.
func (r ProductRepo) UpdateProduct(product entity.Product) error {
	panic("unimplemented")
}

func (r ProductRepo) CreateProduct(ctx context.Context, param entity.Product, userID uint32) (entity.Product, error) {
	var product entity.Product
	query := `
        INSERT INTO "products" (name, sku, category, image_url, notes, price, stock, location, is_available, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id_product
    `

	err := r.dbConnector.DB.QueryRowContext(ctx, query, param.Name,
		param.SKU,
		param.Category,
		param.ImageURL,
		param.Notes,
		param.Price,
		param.Stock,
		param.Location,
		param.IsAvailable,
		param.CreatedAt,
		param.UpdatedAt).Scan(&product.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Product{}, msg.BadRequest("no rows were returned")
		}
		return entity.Product{}, msg.InternalServerError(err.Error())
	}

	return product, nil
}
