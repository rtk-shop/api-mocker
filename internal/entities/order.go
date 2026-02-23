package entities

import gql_gen "rtk/api-mocker/internal/clients/graphql/gen"

type NewOrder struct {
	Name           string                     `fake:"{firstname}"`
	Surname        string                     `fake:"{lastname}"`
	Phone          string                     `fake:"{customerPhoneNumber}"`
	CityName       string                     `fake:"{city}"`
	PostOfficeName string                     `fake:"{streetname}"`
	Supplier       gql_gen.SupplierService    `fake:"{supplierService}"`
	PaymentMethod  gql_gen.OrderPaymentMethod `fake:"{paymentMethod}"`
}

type ProductRangeAmount struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type CreateOrdersRequest struct {
	Quantity           int                `json:"quantity"`
	Products           []string           `json:"products"`
	ProductRangeAmount ProductRangeAmount `json:"productRangeAmount"`
}

type CreateOrdersPayload struct {
	Quantity int                    `json:"quantity"`
	OrdersId []string               `json:"ordersId"`
	Errors   []CreateOrderErrorItem `json:"errors"`
}

type CreateOrderErrorItem struct {
	Message string `json:"message"`
}
