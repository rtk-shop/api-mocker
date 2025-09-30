package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"mime/multipart"
	"net/http"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/entities"
	"strconv"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
)

func RandomSizeName(category gql_gen.CategoryType) string {
	sizes := sizeVariations[category]
	if len(sizes) == 0 {
		return ""
	}

	return sizes[mathRand.Intn(len(sizes))]
}

func (s *service) Create(ctx context.Context, quantity int) (*entities.CreatedProductsPayload, error) {

	s.log.Infof("try to create products, quantity=%d", quantity)

	if _, err := s.loadFiles(); err != nil {
		s.log.Errorw("preload plug files", "error", err)
		return nil, errors.New("failed to preload plug images")
	}

	var wg sync.WaitGroup

	type result struct {
		productID string
		err       error
	}

	resultsCh := make(chan result, quantity)

	for range quantity {

		wg.Add(1)

		go func() {

			defer wg.Done()

			select {
			case <-ctx.Done():
				resultsCh <- result{err: ctx.Err()} // cancled context
				return
			default:
				// continue work
			}

			product, err := s.newMockProduct()
			if err != nil {
				s.log.Errorw("get mock product", "error", err)
				resultsCh <- result{err: err}
				return
			}

			productID, err := s.tempMakeCreateProductMutation(ctx, product)
			if err != nil {
				s.log.Errorw("make createProduct mutation", "error", err)
				resultsCh <- result{err: err}
				return
			}

			resultsCh <- result{productID: productID, err: nil}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var createdProductIDs []string
	var deleteErrors []string

	for res := range resultsCh {
		if res.err != nil {
			deleteErrors = append(deleteErrors, fmt.Sprintf("product_id=%s: %s", res.productID, res.err.Error()))
		} else {
			createdProductIDs = append(createdProductIDs, res.productID)
		}
	}

	s.log.Infow("created products id's", "id", createdProductIDs)

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
		Quantity: len(createdProductIDs),
	}, nil
}

func (s *service) newMockProduct() (entities.NewProduct, error) {

	var maxImages = 2

	var newProduct entities.NewProduct

	err := gofakeit.Struct(&newProduct)
	if err != nil {
		return entities.NewProduct{}, err
	}

	newProduct.Title = newProduct.Title + " " + gofakeit.ProductName()

	// fmt.Printf("%+v\n", newProduct)

	previewFile, ok := s.plugFiles[newProduct.Category]
	if !ok {
		s.log.Errorf("plug image for category not found", "category", newProduct.Category)
		return entities.NewProduct{}, fmt.Errorf("lug image for category %q not found", newProduct.Category)
	}

	newProduct.Preview = previewFile

	otherImages := make([]*entities.ProductImageInput, 0, maxImages)

	for i := range maxImages {

		otherImages = append(otherImages, &entities.ProductImageInput{
			Order: i + 1,
			Image: &previewFile,
		})
	}

	newProduct.Images = otherImages
	newProduct.SizeName = RandomSizeName(newProduct.Category)

	return newProduct, nil

}

// INFO: gql client generator doesn't support nasted graphql.Upload
// issue: https://github.com/Yamashou/gqlgenc/issues/292

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

	files := []*entities.UploadFile{&input.Preview} // индекс 0
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
