package product

import (
	"context"
	"errors"
	"fmt"
	"rtk/api-mocker/internal/entities"
)

func (s *service) Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error) {

	fmt.Println("--->", quantity)

	p, err := s.gql.ProductByID(context.Background(), "7")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(p.Product)

	return nil, errors.New("todo")
}
