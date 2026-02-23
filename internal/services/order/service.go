package order

import (
	"context"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/entities"
	"rtk/api-mocker/internal/utils"
	"rtk/api-mocker/pkg/logger"

	"github.com/brianvoe/gofakeit/v7"
)

type service struct {
	config *config.Config
	log    logger.Logger
	gql    gql_gen.GenGraphQLClient
}

type Service interface {
	Create(ctx context.Context, orders entities.CreateOrdersRequest) (*entities.CreateOrdersPayload, error)
}

type ServiceOptions struct {
	Config    *config.Config
	Logger    logger.Logger
	GqlClient gql_gen.GenGraphQLClient
}

func New(options ServiceOptions) Service {

	gofakeit.AddFuncLookup("customerPhoneNumber", gofakeit.Info{
		Category:    "custom",
		Description: "Generate phone number in RTK format",
		Example:     "099293710",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			operators := []string{"99", "66", "67", "96", "63", "73"}

			s, err := f.Generate("#######")
			if err != nil {
				return nil, err
			}

			return utils.RandomSliceElement(operators) + s, nil
		},
	})

	gofakeit.AddFuncLookup("supplierService", gofakeit.Info{
		Category:    "custom",
		Description: "Generate supplier name in RTK format",
		Example:     "NOVAP",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			suppliers := make([]string, len(gql_gen.AllSupplierService))
			for i, s := range gql_gen.AllSupplierService {
				suppliers[i] = string(s)
			}

			return f.RandomString(suppliers), nil
		},
	})

	gofakeit.AddFuncLookup("paymentMethod", gofakeit.Info{
		Category:    "custom",
		Description: "Generate payment method name in RTK format",
		Example:     "ONLINE",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			methods := make([]string, len(gql_gen.AllOrderPaymentMethod))
			for i, m := range gql_gen.AllOrderPaymentMethod {
				methods[i] = string(m)
			}

			return f.RandomString(methods), nil
		},
	})

	return &service{
		config: options.Config,
		log:    options.Logger,
		gql:    options.GqlClient,
	}
}
