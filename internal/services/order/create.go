package order

import (
	"context"
	"fmt"
	"os"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/entities"
	"rtk/api-mocker/internal/utils"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *service) Create(ctx context.Context, orders entities.CreateOrdersRequest) (*entities.CreateOrdersPayload, error) {

	if orders.Quantity <= 0 {
		return nil, NewInvalidRequestError("quantity should be greater than 0")
	}

	if len(orders.Products) == 0 {
		return nil, NewInvalidRequestError("specify the product IDs for orders")
	}

	amountRange := orders.ProductRangeAmount

	if amountRange == (entities.ProductRangeAmount{}) {
		return nil, NewInvalidRequestError("specify the order's products amount range")
	}

	if amountRange.Min < 1 || amountRange.Max <= amountRange.Min {
		return nil, NewInvalidRequestError("product range values are not correct")
	}

	createdOrderIds := make([]string, 0, orders.Quantity)
	createdOrderErrors := make([]entities.CreateOrderErrorItem, 0, orders.Quantity)

	for range orders.Quantity {

		/* --------------------  Add to cart  ----------------------- */

		for _, productId := range orders.Products {

			amount := utils.RandInt(amountRange.Min, amountRange.Max)

			cartPayload, err := s.gql.AddCartItem(ctx, productId, amount)
			if err != nil {
				if handledError, ok := err.(*clientv2.ErrorResponse); ok {
					fmt.Fprintf(os.Stderr, "handled error: %s\n", handledError.Error())
				} else {
					fmt.Fprintf(os.Stderr, "unhandled error: %s\n", err.Error())
				}

				s.log.Errorw("added product to cart", "error", err.Error())

				createdOrderErrors = append(createdOrderErrors, entities.CreateOrderErrorItem{
					Message: err.Error(),
				})
			}

			cartItem := cartPayload.GetAddCartItem()

			s.log.Infow("added product to cart",
				"product", cartItem.GetProductID(),
				"amount", cartItem.GetQuantity(),
			)
		}

		/* --------------------  Make order  ----------------------- */

		var generatedOrder entities.NewOrder

		err := gofakeit.Struct(&generatedOrder)
		if err != nil {
			return nil, err
		}

		orderPayload, err := s.gql.CreateOrder(ctx, gql_gen.NewOrderInput(generatedOrder))
		if err != nil {
			if handledError, ok := err.(*clientv2.ErrorResponse); ok {
				fmt.Fprintf(os.Stderr, "handled error: %s\n", handledError.Error())
			} else {
				fmt.Fprintf(os.Stderr, "unhandled error: %s\n", err.Error())
			}
			return nil, err
		}

		newOrder := orderPayload.GetCreateOrder()

		s.log.Infow("created order",
			"order_id", newOrder.GetID(),
		)

		createdOrderIds = append(createdOrderIds, newOrder.GetID())
	}

	return &entities.CreateOrdersPayload{
		Quantity: len(createdOrderIds),
		OrdersId: createdOrderIds,
		Errors:   createdOrderErrors,
	}, nil
}
