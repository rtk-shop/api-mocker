package server

import (
	"context"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/server/gen/openapi"
	"rtk/api-mocker/internal/services/product"
	"rtk/api-mocker/pkg/logger"
)

type Server struct {
	config   *config.Config
	log      logger.Logger
	services Services
}

type Services struct {
	Products product.Service
}

func New(config *config.Config, logger logger.Logger, services Services) *Server {
	return &Server{
		config:   config,
		log:      logger,
		services: services,
	}
}

func (s *Server) CreateProducts(ctx context.Context, request openapi.CreateProductsRequestObject) (openapi.CreateProductsResponseObject, error) {
	q := request.Body.Quantity

	if q == 0 {
		return openapi.CreateProducts400JSONResponse{
			Message: "invalid quantity",
		}, nil
	}

	payload, err := s.services.Products.Create(ctx, q)
	if err != nil {
		return openapi.CreateProducts400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	resp := openapi.CreateProductsResponse{
		Quantity: payload.Quantity,
	}

	return openapi.CreateProducts200JSONResponse(resp), nil
}
