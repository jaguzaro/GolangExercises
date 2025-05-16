package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	ctr "TendersGuruAPI/controllers"
	stc "TendersGuruAPI/structs"

	"github.com/go-chi/chi"
)

type APIResponse struct {
	Country string                `json:"country_name"`
	Tenders []stc.TendersGuruData `json:"tenders"`
}

var wg sync.WaitGroup

func main() {
	r := chi.NewRouter()
	r.Get("/", GetTendersGuruHandler)
	fmt.Println("Servidor ejecutandose...")
	http.ListenAndServe(":3000", r)
}

func GetTendersGuruHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	outCountry := make(chan string, 1)
	outTenders := make(chan []stc.TendersGuruData, 1)
	outResponse := make(chan APIResponse, 1)

	wg.Add(2)
	go CallCountryInfo(ctx, cancel, &outCountry)
	go CallTendersGuru(ctx, cancel, &outTenders)

	go func() {
		wg.Wait()
		defer close(outResponse)
		countryName := <-outCountry
		tenders := <-outTenders

		select {
		case <-ctx.Done():
			fmt.Println("Timeout")
		default:
			outResponse <- APIResponse{
				Country: countryName,
				Tenders: tenders,
			}
		}
	}()

	select {
	case <-ctx.Done():
		http.Error(w, "Timeout", http.StatusRequestTimeout)
	case data := <-outResponse:
		fmt.Println(data)
		json.NewEncoder(w).Encode(data)
	}
}

func CallTendersGuru(ctx context.Context, cancel context.CancelFunc, outTenders *chan []stc.TendersGuruData) {
	tenders, err := ctr.GetTendersGuru(ctx)
	defer wg.Done()
	if err != nil {
		cancel()
	}
	*outTenders <- tenders
	close(*outTenders)
}

func CallCountryInfo(ctx context.Context, cancel context.CancelFunc, outCountry *chan string) {
	country, err := ctr.GetCountry(ctx)
	defer wg.Done()
	if err != nil {
		cancel()
	}
	*outCountry <- country
	close(*outCountry)
}
