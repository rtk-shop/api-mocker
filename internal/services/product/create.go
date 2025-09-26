package product

import (
	"context"
	"errors"
	"fmt"
	"rtk/api-mocker/internal/entities"

	"github.com/brianvoe/gofakeit/v7"
)

func (s *service) Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error) {

	s.log.Infof("try to create products, quantity=%d", 11)

	var newProduct entities.NewProduct

	err := gofakeit.Struct(&newProduct)
	if err != nil {
		return nil, err
	}

	// todo
	// newProduct.SizeName
	// newProduct.Preview
	// newProduct.Images

	fmt.Printf("%+v\n", newProduct)

	p, err := s.gql.ProductByID(context.Background(), "7")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(p.Product)

	return nil, errors.New("todo")
}
