package product

import (
	"context"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/entities"
	"rtk/api-mocker/pkg/logger"

	"github.com/brianvoe/gofakeit/v7"
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

	gofakeit.AddFuncLookup("apiGender", gofakeit.Info{
		Category:    "custom",
		Description: "Random product gender name",
		Example:     "FEMALE",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			genders := make([]string, len(gql_gen.AllGender))
			for i, g := range gql_gen.AllGender {
				genders[i] = string(g)
			}

			return f.RandomString(genders), nil
		},
	})

	gofakeit.AddFuncLookup("apiProductCategory", gofakeit.Info{
		Category:    "custom",
		Description: "Random product category name",
		Example:     "SUITCASE",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			categories := make([]string, len(gql_gen.AllCategoryType))
			for i, g := range gql_gen.AllCategoryType {
				categories[i] = string(g)
			}

			return f.RandomString(categories), nil
		},
	})

	return &service{
		config: options.Config,
		log:    options.Logger,
		gql:    options.GqlClient,
	}
}
