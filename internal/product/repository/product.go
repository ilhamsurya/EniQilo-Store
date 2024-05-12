package repository

import (
	"context"
	"database/sql"
	"fmt"
	"projectsphere/eniqlo-store/internal/product/entity"
	"projectsphere/eniqlo-store/pkg/database"
	"projectsphere/eniqlo-store/pkg/protocol/msg"
	"strconv"
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

func (r ProductRepo) GetProduct(ctx context.Context, param entity.GetProductParam) ([]entity.GetProductData, error) {
	query := `
        SELECT p.id_product, p.name, p.category, p.is_available, p.sku, p.stock, p.price, p.created_at, p.notes, p.location, pi.image_url
		FROM "products" p
		JOIN "product_images" pi ON pi.id_product = p.id_product
		WHERE 1=1 
    `

	args := []interface{}{}
	argsCount := 1

	if param.Id != nil {
		query += fmt.Sprintf(" AND p.id_product = $%d", argsCount)
		args = append(args, &param.Id)
		argsCount++
	}

	if param.Name != "" {
		query += " AND p.name ILIKE '%' || $" + strconv.Itoa(argsCount) + " || '%'"
		args = append(args, param.Name)
		argsCount++
	}

	if param.Category != "" {
		query += fmt.Sprintf(" AND p.category = $%d", argsCount)
		args = append(args, param.Category)
		argsCount++
	}

	if param.IsAvailable != nil {
		query += fmt.Sprintf(" AND p.is_available = $%d", argsCount)
		args = append(args, &param.IsAvailable)
		argsCount++
	}

	if param.Sku != "" {
		query += fmt.Sprintf(" AND p.sku = $%d", argsCount)
		args = append(args, param.Sku)
		argsCount++
	}

	if param.InStock != nil {
		if *param.InStock {
			query += " AND p.stock > 0"
		} else {
			query += " AND p.stock <= 0"
		}
	}

	if param.Price != "" || param.CreatedAt != "" {
		orderBy := ""
		if param.Price != "" {
			orderBy += fmt.Sprintf(" p.price %s", param.Price)
		}

		if param.Price != "" && param.CreatedAt != "" {
			orderBy += ","
		}

		if param.CreatedAt != "" {
			orderBy += fmt.Sprintf(" p.created_at %s", param.CreatedAt)
		}

		query += " ORDER BY" + orderBy
	}

	if param.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argsCount)
		args = append(args, *param.Limit)
		argsCount++
	}

	if param.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argsCount)
		args = append(args, &param.Offset)
		argsCount++
	}

	rows, err := r.dbConnector.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return []entity.GetProductData{}, msg.InternalServerError(err.Error())
	}
	defer rows.Close()

	var products = []entity.GetProductData{}
	for rows.Next() {
		var product = entity.GetProductData{}
		var err = rows.Scan(
			&product.ID, &product.Name, &product.Category, &product.IsAvailable, &product.SKU, &product.Stock, &product.Price, &product.CreatedAt, &product.Notes, &product.Location, &product.ImageURL,
		)

		if err != nil {
			return []entity.GetProductData{}, msg.InternalServerError(err.Error())
		}

		products = append(products, product)
	}

	return products, nil
}
