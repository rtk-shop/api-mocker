package product

import (
	"context"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/entities"
	"rtk/api-mocker/pkg/logger"
)

type service struct {
	config *config.Config
	log    logger.Logger
	gql    gql_gen.GenGraphQLClient
}

type Service interface {
	Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error)
}

type ServiceOptions struct {
	Config    *config.Config
	Logger    logger.Logger
	GqlClient gql_gen.GenGraphQLClient
}

func New(options ServiceOptions) Service {

	return &service{
		config: options.Config,
		log:    options.Logger,
		gql:    options.GqlClient,
	}
}
