package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PgRepositoryTestSuite struct {
	suite.Suite
	repo *PgProductRepository
}

func (s *PgRepositoryTestSuite) SetupSuite() {
	url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Printf("Can not connect to database %s", err)
		panic(err)
	}
	s.repo = &PgProductRepository{conn}
}

func (s *PgRepositoryTestSuite) TearDownTest() {
	if _, err := s.repo.cn.Exec(context.Background(), "DELETE FROM products"); err != nil {
		panic(err)
	}
}

func (s *PgRepositoryTestSuite) TestFindAllOrderByExp() {
	err := s.repo.Save(context.Background(), &Product{
		ChatId:         0,
		Name:           "testProduct",
		ProductType:    "FOOD",
		ExpirationDate: time.Now(),
	})
	require.NoError(s.T(), err)

	products, err := s.repo.FindAllOrderByExp(context.Background())
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(products))
}

func TestPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PgRepositoryTestSuite))
}
