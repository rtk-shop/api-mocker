package product

import (
	"context"
	"fmt"
	"rtk/api-mocker/internal/entities"
	"sync"

	"github.com/Yamashou/gqlgenc/clientv2"
)

// todo: use semaphore

func (s *service) Delete(ctx context.Context, productIDs []string) (*entities.DeletedProductsPayload, error) {

	s.log.Infof("try to delete products, quantity=%d", len(productIDs))
	s.log.Infow("product to delete", "id", productIDs)

	var wg sync.WaitGroup

	type result struct {
		productID string
		err       error
	}

	resultsCh := make(chan result)

	for _, id := range productIDs {

		wg.Add(1)

		go func(id string) {

			defer wg.Done()

			payload, err := s.gql.DeleteProduct(ctx, id)
			if err != nil {
				if handledError, ok := err.(*clientv2.ErrorResponse); ok {
					s.log.Errorf("handled error for product_id=%s: %s", id, handledError.Error())
				} else {
					s.log.Errorf("unhandled error for product_id=%s: %s", id, err.Error())
				}

				resultsCh <- result{productID: id, err: err}
			}

			deleteProduct := payload.GetDeleteProduct()

			if deleteProduct != nil {
				resultsCh <- result{productID: deleteProduct.GetID(), err: nil}
			} else {
				resultsCh <- result{productID: id, err: fmt.Errorf("nil response")}
			}

		}(id)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var deletedProductsId []string
	var deleteErrors []string

	for res := range resultsCh {
		if res.err != nil {
			deleteErrors = append(deleteErrors, fmt.Sprintf("product_id=%s: %s", res.productID, res.err.Error()))
		} else {
			deletedProductsId = append(deletedProductsId, res.productID)
		}
	}

	s.log.Infow("deleted products", "id", deletedProductsId)

	if len(deleteErrors) > 0 {
		s.log.Warnw("some products failed to delete", "errors", deleteErrors, "failed_count", len(deleteErrors))
		return nil, fmt.Errorf("some products failed: %s", deleteErrors)
	}

	return &entities.DeletedProductsPayload{
		DeletedQuantity: len(deletedProductsId),
		IDs:             deletedProductsId,
	}, nil
}
