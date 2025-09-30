package entities

import (
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
)

type CreatedProductsPayload struct {
	Quantity int
}

type NewProduct struct {
	Title       string               `fake:"{productname}"`
	SKU         string               `fake:"{productupc}"`
	BasePrice   int                  `fake:"{number:390,8122}"`
	Amount      int                  `fake:"{number:10,30}"`
	Gender      gql_gen.Gender       `fake:"{apiGender}"`
	Category    gql_gen.CategoryType `fake:"{apiProductCategory}"`
	Preview     UploadFile           `fake:"skip"`
	Images      []*ProductImageInput `fake:"skip"`
	Description string               `fake:"{productdescription}"`
	SizeName    string               `fake:"skip"`
	BrandName   string               `fake:"{company}"`
}

type ProductImageInput struct {
	Order int         `json:"order"`
	Image *UploadFile `json:"image"`
}

type UploadFile struct {
	Filename    string
	Data        []byte
	ContentType string
}

type DeletedProductsPayload struct {
	DeletedQuantity int
	IDs             []string
}
