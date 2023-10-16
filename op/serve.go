package op

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Run(port string) {
	http.ListenAndServe(config.Addr, defaultMux)
}

type Controller[ReturnType any, Body any] func(c Ctx[Body]) (ReturnType, error)

type ErrorResponse struct {
	Error string `json:"error"` // human readable error message
}

func HttpHandler[ReturnType any, Body any](controller any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		f, ok := controller.(func(c Ctx[Body]) (ReturnType, error))
		if !ok {
			var c Controller[ReturnType, Body]
			slog.Info("Controller types not ok",
				"type", fmt.Sprintf("%T", controller),
				"should be", fmt.Sprintf("%T", c))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := Ctx[Body]{
			request: r,
		}

		ans, err := f(ctx)

		if err != nil {
			slog.Info("Error in controller", "err", err)
			errResponse := ErrorResponse{
				Error: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}