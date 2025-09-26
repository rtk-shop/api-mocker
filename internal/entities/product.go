package entities

import (
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"

	"github.com/99designs/gqlgen/graphql"
)

type CreatedProductsPayload struct {
	Quantity int
}

type NewProduct struct {
	Title       string                       `fake:"{productname}"`
	SKU         string                       `fake:"{productupc}"`
	BasePrice   int                          `fake:"{number:390,8122}"`
	Amount      int                          `fake:"{number:10,30}"`
	Gender      gql_gen.Gender               `fake:"{apiGender}"`
	Category    gql_gen.CategoryType         `fake:"{apiProductCategory}"`
	Preview     graphql.Upload               `fake:"skip"`
	Images      []*gql_gen.ProductImageInput `fake:"skip"`
	Description string                       `fake:"{productdescription}"`
	SizeName    string                       `fake:"skip"`
	BrandName   string                       `fake:"{company}"`
}
