package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func getInBrasilApi(ctx context.Context, cep string, responseCh chan<- string) {

	select {
	case <-ctx.Done():
		return
	default:
		url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		defer res.Body.Close()

		payload, err := io.ReadAll(res.Body)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		var brasilApiResponse BrasilAPIResponse
		err = json.Unmarshal(payload, &brasilApiResponse)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		fmt.Println(brasilApiResponse)
	}

	responseCh <- "API Utilizada: BRASIL API"
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func getInViaCep(ctx context.Context, cep string, responseCh chan<- string) {
	select {
	case <-ctx.Done():
		return
	default:
		url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		defer res.Body.Close()

		payload, err := io.ReadAll(res.Body)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		var viaCepApiResponse ViaCEPResponse
		err = json.Unmarshal(payload, &viaCepApiResponse)

		if err != nil {
			log.Println(err)
			panic(err)
		}

		fmt.Println(viaCepApiResponse)
	}

	responseCh <- "API Utilizada: VIA CEP"
}

func main() {
	ch := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go getInBrasilApi(ctx, "63111020", ch)
	go getInViaCep(ctx, "63111020", ch)

	fmt.Println(<-ch)

	cancel()
}
