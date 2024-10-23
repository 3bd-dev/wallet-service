package checkapi

import (
	"net/http"
	"os"
	"runtime"

	"github.com/3bd-dev/wallet-service/pkg/database"
	"github.com/3bd-dev/wallet-service/pkg/errs"
	"github.com/3bd-dev/wallet-service/pkg/logger"
	"github.com/3bd-dev/wallet-service/pkg/web"
)

type api struct {
	log *logger.Logger
	db  database.IDatabase
}

func newapi(db database.IDatabase, log *logger.Logger) *api {
	return &api{log: log, db: db}
}

// readiness checks if the database is ready and if not will return a 500 status.
func (a *api) readiness(w http.ResponseWriter, r *http.Request) {
	if err := a.db.Ping(); err != nil {
		a.log.Info(r.Context(), "readiness failure", "ERROR", err)
		errs.New(errs.Internal, err)
		return
	}
	web.RenderNoContent(w)
}

// liveness returns simple status info if the service is alive.
func (a *api) liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := Info{
		Status:     "up",
		Host:       host,
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// This handler provides a free timer loop.
	web.RenderOk(w, info)
}
