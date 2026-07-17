package api

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/utils"
)

// Health reports whether the service is up and can reach its database. Uptime
// checks and the container healthcheck hit this.
func Health(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			utils.JsonResponse(w, false, http.StatusServiceUnavailable, "database unreachable", nil)
			return
		}
		utils.JsonResponse(w, true, http.StatusOK, "ok", nil)
	}
}
