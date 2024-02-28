package main

import (
	"net/http"
	"os"

	"github.com/danielzinhors/cloudrun_go/internal/handlers"
	"github.com/danielzinhors/cloudrun_go/internal/services"
	"github.com/danielzinhors/cloudrun_go/internal/usecases"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	getTempHandler := handlers.NewGetTempHandler(
		usecases.NewGetTempUseCase(
			services.NewViaCepService(),
			services.NewWeatherApiService(),
		),
	)

	r.Get("/sol", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/", getTempHandler.Handle)
	portServer := ":" + os.Getenv("PORT_API")

	err := http.ListenAndServe(portServer, r)
	if err != nil {
		panic(err)
	}
}
