package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/generated/openapi"
	"rtk/api-mocker/pkg/logger"
	"runtime"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	flagLocal  = "local"
	flagDocker = "docker"
)

func init() {

	info, _ := debug.ReadBuildInfo()
	fmt.Println("Go version:", info.GoVersion, runtime.GOARCH)

	// without default value
	env := flag.String("env", "", "specify .env filename for flag")
	flag.Parse()
	// for dynamic Load(".env." + *env);

	if *env == flagLocal {
		if err := godotenv.Load(".env.local", ".env"); err != nil {
			log.Fatal(err)
		}
		log.Printf("⚙️ loaded .env.%s, .env\n", *env)
	}

	if *env == flagDocker {
		log.Printf("🐳 app runs in %q mod", flagDocker)
	}

}

func main() {

	config := config.New()

	zap := logger.New(config)

	zap.Info("info message")
	// zap.Debug("debug message")
	// zap.Warn("warn message")
	// zap.Error("error message")

	// r := chi.NewRouter()

	// gqlClient := gql_gen.NewClient(http.DefaultClient, config.ApiURL, nil,
	// 	func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res any, next clientv2.RequestInterceptorFunc) error {
	// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ApiToken))

	// 		return next(ctx, req, gqlInfo, res)
	// 	})

	// p, err := gqlClient.ProductByID(context.Background(), "7")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(p.Product)

	r := chi.NewRouter()

	strictHandler := openapi.NewStrictHandler(&Server{}, nil)

	handler := openapi.HandlerFromMux(strictHandler, r)

	if err := http.ListenAndServe(":"+config.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	// app := app.NewApp(config, logger, db, r2store)
	// app.Run()
}

type Server struct{}

func (s *Server) CreateProducts(ctx context.Context, request openapi.CreateProductsRequestObject) (openapi.CreateProductsResponseObject, error) {
	q := request.Body.Quantity

	if q == 0 {
		return openapi.CreateProducts400JSONResponse{
			Message: "invalid quantity",
		}, nil
	}

	resp := openapi.CreateProductsResponse{
		Quantity: q,
	}

	return openapi.CreateProducts200JSONResponse(resp), nil
}
