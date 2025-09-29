package product

import (
	"context"
	"io"
	"net/http"
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
	Delete(ctx context.Context, productIDs []string) (*entities.DeletedProductsPayload, error)
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

func fetchFile(url, filename string) (*entities.UploadFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &entities.UploadFile{
		Filename:    filename,
		Data:        data,
		ContentType: http.DetectContentType(data),
	}, nil
}

// func downloadAsUpload(url, filename string) (graphql.Upload, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return graphql.Upload{}, err
// 	}

// 	defer resp.Body.Close()

// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return graphql.Upload{}, err
// 	}

// 	upload := graphql.Upload{
// 		File:        bytes.NewReader(data),
// 		Filename:    filename,
// 		Size:        int64(len(data)),
// 		ContentType: http.DetectContentType(data),
// 	}

// 	return upload, nil
// }
