package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type ProductRepository interface {
	// Save new product
	Save(ctx context.Context, p *Product) error
	// FindAllOrderByExp Return all not expired products order by expiration date
	FindAllOrderByExp(ctx context.Context, isExpired bool) ([]Product, error)
	// Update existing products by id
	Update(ctx context.Context, id int64, p *Product) error
	// FindByExpiredTime return products that expired in given duration
	FindByExpiredTime(ctx context.Context, interval time.Duration, productType string) ([]Product, error)
	// Delete product by id
	Delete(ctx context.Context, idForDel int64) error
	// SetExpired set all expired products
	SetExpired(ctx context.Context) error
}

type Product struct {
	Id             int64
	ChatId         int64
	Name           string
	ProductType    string
	ExpirationDate time.Time
}

type RepositoryConfig struct {
	Url string
}

type PgProductRepository struct {
	cn *pgx.Conn
}

func NewPgProductRepository(cfg *RepositoryConfig) (ProductRepository, error) {
	conn, err := pgx.Connect(context.Background(), cfg.Url)
	if err != nil {
		return nil, err
	}
	return &PgProductRepository{conn}, nil
}

func (pg *PgProductRepository) Save(ctx context.Context, p *Product) error {
	_, err := pg.cn.Exec(ctx, "INSERT INTO products(chat_id, name, type, expiration_date) VALUES($1,$2,$3,$4)",
		p.ChatId, p.Name, p.ProductType, p.ExpirationDate)
	return err
}

func (pg *PgProductRepository) FindAllOrderByExp(ctx context.Context, isExpired bool) ([]Product, error) {
	rows, err := pg.cn.Query(ctx,
		"SELECT id, chat_id, name, type, expiration_date FROM products\n WHERE is_expired = $1 ORDER BY expiration_date DESC", isExpired)
	if err != nil {
		return nil, err
	}
	return mapRowsToProducts(rows)
}

func (pg *PgProductRepository) Update(ctx context.Context, id int64, p *Product) error {
	_, err := pg.cn.Exec(ctx, "UPDATE products SET chat_id = $2, name = $3, type = $4, expiration_date = $5, is_expired = now() > $5  WHERE id = $1",
		id, p.ChatId, p.Name, p.ProductType, p.ExpirationDate)
	return err
}

func (pg *PgProductRepository) FindByExpiredTime(ctx context.Context, interval time.Duration, productType string) ([]Product, error) {
	r, err := pg.cn.Query(ctx, "SELECT id, chat_id, name, type, expiration_date FROM products WHERE type = $1 AND expiration_date <= now() + $2", productType, interval)
	if err != nil {
		return nil, err
	}
	return mapRowsToProducts(r)
}

func (pg *PgProductRepository) Delete(ctx context.Context, idForDel int64) error {
	_, err := pg.cn.Exec(ctx, "DELETE from products WHERE id = $1", idForDel)
	return err
}

func (pg *PgProductRepository) SetExpired(ctx context.Context) error {
	_, err := pg.cn.Exec(ctx, "UPDATE products SET is_expired = TRUE WHERE now() > expiration_date ")
	return err
}

func mapRowsToProducts(r pgx.Rows) ([]Product, error) {
	if err := r.Err(); err != nil {
		return nil, err
	}
	products := make([]Product, 0)
	for r.Next() {
		p := Product{}
		err := r.Scan(&p.Id, &p.ChatId, &p.Name, &p.ProductType, &p.ExpirationDate)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
