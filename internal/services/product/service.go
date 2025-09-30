package product

import (
	"context"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/entities"
	"rtk/api-mocker/pkg/logger"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
)

type service struct {
	config    *config.Config
	log       logger.Logger
	gql       gql_gen.GenGraphQLClient
	plugFiles map[gql_gen.CategoryType]entities.UploadFile
	once      sync.Once
}

type Service interface {
	Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error)
	Delete(ctx context.Context, productIDs []string) (*entities.DeletedProductsPayload, error)
}

type ServiceOptions struct {
	Config    *config.Config
	Logger    logger.Logger
	GqlClient gql_gen.GenGraphQLClient
}

var plugImagesURL = map[gql_gen.CategoryType]string{
	gql_gen.CategoryTypeBackpack: "https://s3.rtkstore.org/plug/backpack.jpg",
	gql_gen.CategoryTypeBag:      "https://s3.rtkstore.org/plug/bag.jpg",
	gql_gen.CategoryTypeOther:    "https://s3.rtkstore.org/plug/other.jpg",
	gql_gen.CategoryTypeSuitcase: "https://s3.rtkstore.org/plug/suitcase.jpg",
}

var sizeVariations = map[gql_gen.CategoryType][]string{
	gql_gen.CategoryTypeSuitcase: {"S", "M", "L"},
	gql_gen.CategoryTypeBackpack: {"S", "M"},
	gql_gen.CategoryTypeBag:      {"S", "M"},
	gql_gen.CategoryTypeOther:    {"none"},
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

	gofakeit.AddFuncLookup("productname10", gofakeit.Info{
		Category:    "custom",
		Description: "Product name with min length 10",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {

			for {
				name := gofakeit.ProductName()
				if len(name) >= 10 {
					return name, nil
				}
			}
		},
	})

	return &service{
		config: options.Config,
		log:    options.Logger,
		gql:    options.GqlClient,
	}
}
