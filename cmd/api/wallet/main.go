package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/3bd-dev/wallet-service/config"
	"github.com/3bd-dev/wallet-service/internal/handlers/checkapi"
	"github.com/3bd-dev/wallet-service/internal/handlers/walletapi"
	"github.com/3bd-dev/wallet-service/internal/models"
	"github.com/3bd-dev/wallet-service/internal/payment"
	"github.com/3bd-dev/wallet-service/internal/payment/gateways/gatewaya"
	"github.com/3bd-dev/wallet-service/internal/payment/gateways/gatewayb"
	"github.com/3bd-dev/wallet-service/internal/repos/postgres"
	"github.com/3bd-dev/wallet-service/internal/services/wallet"
	"github.com/3bd-dev/wallet-service/internal/web/mid"
	"github.com/3bd-dev/wallet-service/pkg/database"
	"github.com/3bd-dev/wallet-service/pkg/database/psql"
	"github.com/3bd-dev/wallet-service/pkg/logger"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

var build = "develop"

func main() {
	log := logger.New(os.Stdout, logger.LevelInfo, "WALLET")

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// Configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	// Database setup
	db, err := psql.Open(psql.Config{
		URL:          cfg.Database.URL,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize db: %w", err)
	}
	database.DBConn = db
	defer db.Close()

	// Repository setup
	walletRepo := postgres.NewWalletRepo(db)
	transactionRepo := postgres.NewTransactionRepo(db)

	// Payment gateway setup
	paymentGateways := map[models.PaymentGateway]payment.PaymentGateway{
		models.PaymentGatewayA: gatewaya.New(cfg.PaymentGatewayConfig.GatewayA),
		models.PaymentGatewayB: gatewayb.New(cfg.PaymentGatewayConfig.GatewayB),
	}

	// Payment handler setup
	paymentHandler := payment.New(paymentGateways)

	// Wallet service setup
	walletService := wallet.NewService(log, walletRepo, transactionRepo, paymentHandler, cfg.PaymentGatewayConfig.CallbackPattern)
	walletService.Start(ctx)

	// HTTP server setup
	httpmux := mux.NewRouter()

	checkapi.Routes(httpmux, checkapi.Config{
		DB:  db,
		Log: log,
	})

	httpmux.Use(mid.Logger(log))

	walletapi.Routes(httpmux, walletapi.Config{
		Service: walletService,
	})

	// swagger setup
	if cfg.Service.Environment == "development" {
		serveSwagger(httpmux)
	}

	httpserver := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      httpmux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", httpserver.Addr)
		serverErrors <- httpserver.ListenAndServe()
	}()

	// Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := httpserver.Shutdown(ctx); err != nil {
			httpserver.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func serveSwagger(r *mux.Router) {
	r.HandleFunc("/swagger/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		yamlPath, _ := filepath.Abs(filepath.Join("docs", "swagger.yaml"))
		http.ServeFile(w, r, yamlPath)
	}).Methods(http.MethodGet)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}
