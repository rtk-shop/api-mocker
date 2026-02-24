package server

import (
	"context"
	"errors"
	"rtk/api-mocker/internal/server/gen/openapi"
	"rtk/api-mocker/internal/services/order"
	"time"
)

func (s *Server) CreateOrders(ctx context.Context, request openapi.CreateOrdersRequestObject) (openapi.CreateOrdersResponseObject, error) {

	start := time.Now()

	payload, err := s.services.Orders.Create(ctx, *request.Body)
	if err != nil {
		var inv *order.InvalidRequestError
		if errors.As(err, &inv) {

			return openapi.CreateOrders422JSONResponse{
				Reason:  inv.Reason,
				Message: inv.Message,
			}, nil
		}

		return openapi.CreateOrders400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	s.log.Infof("%q execution duration time=%s\n", "create-orders", time.Since(start))

	errors := make([]openapi.CreateOrderErrorItem, 0, len(payload.Errors))

	for _, respErr := range payload.Errors {
		errors = append(errors, openapi.CreateOrderErrorItem{
			Message: respErr.Message,
		})
	}

	resp := openapi.CreateOrdersResponse{
		Quantity: payload.Quantity,
		OrdersId: payload.OrdersId,
		Errors:   errors,
	}

	return openapi.CreateOrders200JSONResponse(resp), nil
}
