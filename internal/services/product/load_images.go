package product

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	gql_gen "rtk/api-mocker/internal/clients/graphql/gen"
	"rtk/api-mocker/internal/entities"
	"sync"
)

type fileResult struct {
	key  gql_gen.CategoryType
	file entities.UploadFile
	err  error
}

func (s *service) loadFiles() (map[gql_gen.CategoryType]entities.UploadFile, error) {

	var onceErr error

	s.once.Do(func() {

		s.log.Info("preload plug images")

		s.plugFiles = make(map[gql_gen.CategoryType]entities.UploadFile, len(plugImagesURL))
		results := make(chan fileResult, len(plugImagesURL))

		var wg sync.WaitGroup

		for key, url := range plugImagesURL {
			wg.Add(1)

			// no params, cause go 1.24+
			go func() {
				defer wg.Done()

				file, err := s.fetchFile(url)
				if err != nil {
					results <- fileResult{key: key, err: fmt.Errorf("failed to fetch %s: %w", url, err)}
				}

				results <- fileResult{key: key, file: file}
			}()
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		for res := range results {
			if res.err != nil {
				onceErr = res.err
				return
			}
			s.plugFiles[res.key] = res.file
		}
	})

	if onceErr != nil {
		return nil, onceErr
	}

	return s.plugFiles, nil
}

func (s *service) fetchFile(url string) (entities.UploadFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return entities.UploadFile{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.UploadFile{}, err
	}

	return entities.UploadFile{
		Filename:    rand.Text() + ".jpg",
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
