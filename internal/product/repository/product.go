package repository

import (
	"context"
	"database/sql"
	"projectsphere/eniqlo-store/internal/product/entity"
	"projectsphere/eniqlo-store/pkg/database"
	"projectsphere/eniqlo-store/pkg/protocol/msg"
)

type ProductRepo struct {
	dbConnector database.PostgresConnector
}

func NewProductRepo(dbConnector database.PostgresConnector) ProductRepo {
	return ProductRepo{
		dbConnector: dbConnector,
	}
}

func (r ProductRepo) UpdateProduct(product entity.Product) error {
	query := `
        UPDATE "products"
        SET name = $1, sku = $2, category = $3, image_url = $4, notes = $5, price = $6, stock = $7, location = $8, is_available = $9, updated_at = $10
        WHERE id_product = $11
    `

	_, err := r.dbConnector.DB.Exec(query,
		product.Name,
		product.SKU,
		product.Category,
		product.ImageURL,
		product.Notes,
		product.Price,
		product.Stock,
		product.Location,
		product.IsAvailable,
		product.UpdatedAt,
		product.ID)

	if err != nil {
		return msg.InternalServerError(err.Error())
	}

	return nil
}

func (r ProductRepo) DeleteProduct(id string, userId uint32) error {
	query := `
        DELETE FROM "products"
        WHERE id_product = $1 AND user_id = $2
    `

	result, err := r.dbConnector.DB.Exec(query, id, userId)
	if err != nil {
		return msg.InternalServerError(err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return msg.NotFound("product not found")
	}

	return nil
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
