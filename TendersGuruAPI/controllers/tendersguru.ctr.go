package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	stc "TendersGuruAPI/structs"
)

func GetTendersGuru(ctx context.Context) ([]stc.TendersGuruData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://tenders.guru/api/hu/tenders", nil)
	if err != nil {
		fmt.Println("Error al crear la peticion")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la peticion")
	}
	defer resp.Body.Close()

	var tendersGuruResponse stc.TendersGuruResponse
	if err := json.NewDecoder(resp.Body).Decode(&tendersGuruResponse); err != nil {
		fmt.Println("Error al obtener resultados")
	}

	data := tendersGuruResponse.Data
	if len(tendersGuruResponse.Data) > 15 {
		data = tendersGuruResponse.Data[:15]
	}

	return data, nil
}
