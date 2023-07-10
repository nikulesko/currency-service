package handlers

import (
	"currency-service/internal/storage"
	"currency-service/internal/storage/sqlite"

	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

func New(log *slog.Logger, storage *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fnPath = "http-server.handlers.New"

		log = log.With(slog.String("func path", fnPath))

		date := r.URL.Query().Get("date")
		if date != "" {
			res, err := storage.GetCurrencyByDate(date)
			if err != nil {
				log.Error("Date parameter does not found", err)
		
				render.JSON(w, r, http.Response{
					StatusCode: 400,
					Status: err.Error(),
				})
		
				return
			}

			render.JSON(w, r, res.Stringify())
	
			return
		} else {
			rates, err := RESTCall(log)
			if err != nil {
				log.Error("Cannot call currency service", err)
		
				render.JSON(w, r, http.Response{
					StatusCode: 500,
					Status: err.Error(),
				})
		
				return
			}
			
			err = storage.SaveCurrency(rates)
			if err != nil {
				log.Error("Cannot save currency to the database", err)
			}
			render.JSON(w, r, rates.Stringify())
	
			return
		}
	}
}

func RESTCall(log *slog.Logger) (*storage.LatestRates, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", "https://api.exchangerate.host/latest?base=USD&symbols=EUR,JPY,UAH", nil)
	if err != nil {
		log.Error("Request build error", err)
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Error("Response error", err)
		return nil, err
	}

	var result storage.RawRates
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Clean(), err
}