package checkapi

import (
	"net/http"

	"github.com/3bd-dev/wallet-service/pkg/database"
	"github.com/3bd-dev/wallet-service/pkg/logger"
	"github.com/gorilla/mux"
)

type Config struct {
	Log *logger.Logger
	DB  database.IDatabase
}

// Routes adds specific routes for this group.
func Routes(router *mux.Router, cfg Config) {
	api := newapi(cfg.DB, cfg.Log)

	router.HandleFunc("/readiness", api.readiness).Methods(http.MethodGet)
	router.HandleFunc("/liveness", api.liveness).Methods(http.MethodGet)
}
