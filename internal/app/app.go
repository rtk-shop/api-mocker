package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rtk/api-mocker/internal/config"
	"rtk/api-mocker/internal/generated/openapi"
	"rtk/api-mocker/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

type App struct {
	config      *config.Config
	log         logger.Logger
	httpHandler http.Handler
}

func New(config *config.Config, logger logger.Logger) App {

	r := chi.NewRouter()

	// gqlClient := gql_gen.NewClient(http.DefaultClient, config.ApiURL, nil,
	// 	func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res any, next clientv2.RequestInterceptorFunc) error {
	// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ApiToken))

	// 		return next(ctx, req, gqlInfo, res)
	// 	})

	// p, err := gqlClient.ProductByID(context.Background(), "7")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return App{}
	// }

	// fmt.Println(p.Product)

	strictHandler := openapi.NewStrictHandler(&Server{}, nil)

	handler := openapi.HandlerFromMux(strictHandler, r)

	return App{
		config:      config,
		log:         logger,
		httpHandler: handler,
	}
}

func (a *App) Run() {

	server := &http.Server{
		Addr:              ":" + a.config.Port,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      10 * time.Second,
		Handler:           http.TimeoutHandler(a.httpHandler, 15*time.Second, "request timeout expired"),
	}

	a.log.Info(fmt.Sprintf("http server is running http://localhost%s", server.Addr))

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("shutdown reason: %s \n", err)
	}

	log.Print("✅ server shutdown gracefully")
}

// ====================================================================================================

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
