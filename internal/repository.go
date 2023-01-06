package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type ProductRepository interface {
	Save(ctx context.Context, p *Product) error
	FindAllOrderByExp(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id int64, p *Product) error
	FindByExpiredTime(ctx context.Context, time time.Time) ([]Product, error)
	Delete(ctx context.Context, id int64) error
	SetExpired(ctx context.Context) error
	FindAllExpired(ctx context.Context) ([]Product, error)
}

type Product struct {
	Id             int64
	ChatId         int64
	Name           string
	ProductType    string
	ExpirationDate time.Time
}

type PgProductRepository struct {
	cn *pgx.Conn
}

func (pg *PgProductRepository) Save(ctx context.Context, p *Product) error {
	_, err := pg.cn.Exec(ctx, "INSERT INTO products(chat_id, name, type, expiration_date) VALUES($1,$2,$3,$4)",
		p.ChatId, p.Name, p.ProductType, p.ExpirationDate)
	return err
}

func (pg *PgProductRepository) FindAllOrderByExp(ctx context.Context) ([]Product, error) {
	rows, err := pg.cn.Query(ctx,
		"SELECT id, chat_id, name, type, expiration_date FROM products\n WHERE is_expired = FALSE ORDER BY expiration_date DESC")
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	products := make([]Product, 0)
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.Id, &p.ChatId, &p.Name, &p.ProductType, &p.ExpirationDate)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (pg *PgProductRepository) Update(ctx context.Context, id int64, p *Product) error {
	_, err := pg.cn.Exec(ctx, "UPDATE products SET chat_id = $2, name = $3, type = $4, expiration_date = $5, is_expired = now() > $5  WHERE id = $1",
		id, p.ChatId, p.Name, p.ProductType, p.ExpirationDate)
	return err
}

func (pg *PgProductRepository) FindByExpiredTime(ctx context.Context, time time.Time) ([]Product, error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PgProductRepository) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PgProductRepository) SetExpired(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PgProductRepository) FindAllExpired(ctx context.Context) ([]Product, error) {
	//TODO implement me
	panic("implement me")
}
