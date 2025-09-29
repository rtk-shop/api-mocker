package product

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	mathRand "math/rand"
	"mime/multipart"
	"net/http"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/entities"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
)

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

func RandomSizeName(category gql_gen.CategoryType) string {
	sizes := sizeVariations[category]
	if len(sizes) == 0 {
		return ""
	}

	return sizes[mathRand.Intn(len(sizes))]
}

func (s *service) Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error) {

	s.log.Infof("try to create products, quantity=%d", quantity)

	// INFO: gql client generator doesn't support nasted graphql.Upload
	// issue: https://github.com/Yamashou/gqlgenc/issues/292

	// previewUpload, err := downloadAsUpload(plugImagesURL[newProduct.Category], rand.Text()+".jpg")
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Printf("Upload готов: filename=%s, size=%d, contentType=%s\n",
	// 	previewUpload.Filename, previewUpload.Size, previewUpload.ContentType)

	// imagesCount := 2

	// otherImages := make([]*gql_gen.ProductImageInput, 0, imagesCount)

	// for i := range imagesCount {
	// 	// fmt.Println(previewUpload)

	// 	img, err := downloadAsUpload(plugImagesURL[newProduct.Category], rand.Text()+".jpg")
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	otherImages = append(otherImages, &gql_gen.ProductImageInput{
	// 		Order: i + 1,
	// 		Image: img,
	// 	})
	// }

	// fmt.Println("otherImages", len(otherImages), cap(otherImages))
	// for _, img := range otherImages {
	// 	fmt.Printf("order=%d, filename=%s, size=%d\n", img.Order, img.Image.Filename, img.Image.Size)
	// }

	var createdIDs []string

	for range quantity {

		product, err := s.newMockProduct()
		if err != nil {
			s.log.Errorw("get mock product", "error", err)
			return nil, err
		}

		productID, err := s.tempMakeCreateProductMutation(ctx, product)
		if err != nil {
			s.log.Errorw("make createProduct mutation", "error", err)
			return nil, err
		}

		createdIDs = append(createdIDs, productID)
	}

	s.log.Infow("created products id's", "id", createdIDs)

	// createProduct, err := s.gql.CreateProduct(ctx,
	// 	newProduct.Title,
	// 	newProduct.SKU,
	// 	float64(newProduct.BasePrice),
	// 	newProduct.Amount,
	// 	newProduct.Gender,
	// 	newProduct.Category,
	// 	previewUpload,
	// 	otherImages,
	// 	newProduct.Description,
	// 	"XL",
	// 	newProduct.BrandName,
	// )
	// if err != nil {
	// 	if handledError, ok := err.(*clientv2.ErrorResponse); ok {
	// 		fmt.Fprintf(os.Stderr, "handled error: %s\n", handledError.Error())
	// 	} else {
	// 		fmt.Fprintf(os.Stderr, "unhandled error: %s\n", err.Error())
	// 	}

	// 	return nil, err
	// }

	// fmt.Println("-->", createProduct)

	return &entities.CreatedProductsPayload{
		Quantity: len(createdIDs),
	}, nil
}

func (s *service) newMockProduct() (entities.NewProduct, error) {

	var newProduct entities.NewProduct

	err := gofakeit.Struct(&newProduct)
	if err != nil {
		return entities.NewProduct{}, err
	}

	// fmt.Printf("%+v\n", newProduct)

	previewFile, err := fetchFile(plugImagesURL[newProduct.Category], rand.Text()+".jpg")
	if err != nil {
		return entities.NewProduct{}, err
	}

	newProduct.Preview = previewFile

	otherImages := make([]*entities.ProductImageInput, 0, 2)

	for i := range 2 {

		img, err := fetchFile(plugImagesURL[newProduct.Category], rand.Text()+".jpg")
		if err != nil {
			return entities.NewProduct{}, err
		}

		otherImages = append(otherImages, &entities.ProductImageInput{
			Order: i + 1,
			Image: img,
		})
	}

	newProduct.Images = otherImages
	newProduct.SizeName = RandomSizeName(newProduct.Category)

	return newProduct, nil

}

func (s *service) tempMakeCreateProductMutation(ctx context.Context, input entities.NewProduct) (string, error) {

	query := gql_gen.CreateProductDocument

	variables := map[string]any{
		"title":       input.Title,
		"sku":         input.SKU,
		"basePrice":   input.BasePrice,
		"amount":      input.Amount,
		"gender":      input.Gender,
		"category":    input.Category,
		"preview":     nil, // Upload
		"images":      make([]map[string]any, len(input.Images)),
		"description": input.Description,
		"sizeName":    input.SizeName,
		"brandName":   input.BrandName,
	}

	for i, img := range input.Images {
		variables["images"].([]map[string]any)[i] = map[string]any{
			"order": img.Order,
			"image": nil,
		}
	}

	files := []*entities.UploadFile{input.Preview} // индекс 0
	mapData := map[string][]string{
		"0": {"variables.preview"}, // путь к preview
	}

	for i, img := range input.Images {
		files = append(files, img.Image)
		mapData[fmt.Sprintf("%d", i+1)] = []string{fmt.Sprintf("variables.images.%d.image", i)}
	}

	// create multipart body
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	opsJSON, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal operations: %w", err)
	}

	w.WriteField("operations", string(opsJSON))

	// fmt.Println(string(opsJSON))

	mapJSON, err := json.Marshal(mapData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal map data: %w", err)
	}

	w.WriteField("map", string(mapJSON))

	for idx, file := range files {
		if file == nil {
			continue
		}

		part, err := w.CreateFormFile(strconv.Itoa(idx), file.Filename)
		if err != nil {
			return "", fmt.Errorf("failed to create form-file for index %d: %w", idx, err)
		}

		if _, err := part.Write(file.Data); err != nil {
			return "", fmt.Errorf("failed to write file data for index %d: %w", idx, err)
		}
	}

	w.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ApiURL, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.config.ApiToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server error: status %d, body: %s", resp.StatusCode, payload)
	}

	// fmt.Println("Server response:", string(payload))

	gqlResp, err := parseGraphQLResponse(payload)
	if err != nil {
		return "", fmt.Errorf("server error: parse gql response, payload: %s", payload)
	}

	s.log.Info("✅ product created: ", string(gqlResp.Data))

	// {"createProduct":{"id":"234","title":"Robust Appliance Dash","currentPrice":6754,"basePrice":6754}}

	var createProductResp GraphQLCreateProductResponse
	if err := json.Unmarshal(gqlResp.Data, &createProductResp); err != nil {
		return "", fmt.Errorf("server error: parse createProduct raw response, err: %w", err)
	}

	return createProductResp.CreateProduct.ID, nil
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type GraphQLCreateProductResponse struct {
	CreateProduct struct {
		ID           string  `json:"id"`
		Title        string  `json:"title"`
		CurrentPrice float64 `json:"currentPrice"`
		BasePrice    float64 `json:"basePrice"`
	} `json:"createProduct"`
}

func parseGraphQLResponse(body []byte) (*GraphQLResponse, error) {

	var resp GraphQLResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	// with "errors"
	if len(resp.Errors) > 0 {
		return &resp, fmt.Errorf("graphql error: %s (path=%v)",
			resp.Errors[0].Message, resp.Errors[0].Path)
	}

	// with "data"
	if len(resp.Data) > 0 && string(resp.Data) != "null" {
		return &resp, nil
	}

	// no data or errors — unexpected
	return &resp, fmt.Errorf("unexpected graphql response: %s", body)
}
