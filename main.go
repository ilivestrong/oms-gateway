package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	internal "github.com/ilivestrong/oms-gateway/internal"
	"github.com/ilivestrong/oms-gateway/internal/auth"
	"github.com/ilivestrong/oms-gateway/internal/gatewayservice"
	"github.com/ilivestrong/oms-gateway/internal/middlewares"
	omspb "github.com/ilivestrong/oms-gateway/internal/protos"
	env "github.com/joho/godotenv"
	"github.com/justinas/alice"
)

const (
	envFile = ".env"
	version = "v1.0.0"
)

var (
	loadEnv                       = env.Load
	EncodingTypeJSON       string = "json"
	ErrInvalidTokenRequest        = "email missing in the request"
)

func main() {
	appLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := loadEnv(envFile)
	if err != nil {
		log.Fatal("failed to load .env")
	}
	gatwayAddr, exist := os.LookupEnv("LISTEN_ADDRESS_HTTP")
	if !exist {
		fmt.Println("no port specified, defaulting to 5015")
		gatwayAddr = "5015"
	}
	productSvcAddress, exist := os.LookupEnv("LISTEN_ADDRESS_PRODUCT")
	if !exist {
		log.Fatal("invalid product service address")
	}
	OrderSvcAddress, exist := os.LookupEnv("LISTEN_ADDRESS_ORDER")
	if !exist {
		log.Fatal("invalid order service address")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := &internal.Options{
		ListenAddressHTTPPort:       gatwayAddr,
		OrderServiceListenAddress:   OrderSvcAddress,
		ProductServiceListenAddress: productSvcAddress,
	}

	appLogger.Info("oms-gatway", "version", version)
	runGatewayServer(ctx, options, appLogger)
}

func runGatewayServer(ctx context.Context, opts *internal.Options, logger *slog.Logger) {
	svc, err := internal.New(opts)
	if err != nil {
		log.Fatal("gateway service failed to start.")
	}

	orderSvcClient := omspb.NewOrderServiceClient(svc.OrderSvcClientConn)
	productSvcClient := omspb.NewProductServiceClient(svc.ProductSvcClientConn)

	gatewaySvc := gatewayservice.New(productSvcClient, orderSvcClient)

	mux := runtime.NewServeMux()
	muxWithMiddlewares := bindMiddlewaresToMux(mux, middlewares.Authorize, middlewares.RateLimitMiddleware(10, time.Second))
	muxWithMiddlewares.HandleFunc("/login", authHandler(logger))

	if err := omspb.RegisterGatewayServiceHandlerServer(ctx, mux, gatewaySvc); err != nil {
		log.Fatalf("faild to register: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(":"+opts.ListenAddressHTTPPort, muxWithMiddlewares); err != nil {
			log.Fatalf("Failed to start server:: http.ListenAndServe(): %v", err)
		}
	}()
	logger.Info("server listening at:", "port", opts.ListenAddressHTTPPort)

	shutdownOnSignal(svc)
}

func bindMiddlewaresToMux(mux *runtime.ServeMux, mws ...alice.Constructor) *http.ServeMux {
	muxWithMiddlewares := http.NewServeMux()
	muxWithMiddlewares.Handle("/", alice.New(mws...).Then(mux))
	return muxWithMiddlewares
}

func authHandler(logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenRequest, err := getTokenRequest(r)
		if err != nil {
			sendResponse(w, []byte(err.Error()), "", http.StatusBadRequest)
			return
		}

		if strings.Trim(tokenRequest.Email, " ") == "" {
			sendResponse(w, []byte(ErrInvalidTokenRequest), "", http.StatusBadRequest)
			return
		}

		token, err := auth.GenerateAccessToken(*tokenRequest)
		if err != nil {
			log.Fatal(err)
		}
		tokenBytes, err := json.Marshal(token)
		if err != nil {
			logger.Error("authHandler:", "err", fmt.Sprintf("failed to create token: %v", err))
			sendResponse(w, []byte("failed to create token"), "", http.StatusInternalServerError)
			return
		}
		sendResponse(w, tokenBytes, EncodingTypeJSON, http.StatusOK)
	}
}

func getTokenRequest(req *http.Request) (*auth.TokenRequest, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("failed to read login request")
	}
	defer req.Body.Close()

	var tokenRequest auth.TokenRequest
	err = json.Unmarshal(body, &tokenRequest)
	if err != nil {
		return nil, errors.New("failed parse request data")
	}
	return &tokenRequest, nil
}

func sendResponse(w http.ResponseWriter, bodyBytes []byte, encoding string, status int) {
	if encoding == EncodingTypeJSON {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	w.Write(bodyBytes)
}

func waitForShutdownSignal() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c

	return sig.String()
}

func shutdownOnSignal(svc *internal.Service) {
	signalName := waitForShutdownSignal()
	fmt.Printf("recieved signal: %s starting shutdown...", signalName)

	if svc.OrderSvcClientConn != nil {
		svc.OrderSvcClientConn.Close()
	}

	if svc.ProductSvcClientConn != nil {
		svc.ProductSvcClientConn.Close()
	}
}
